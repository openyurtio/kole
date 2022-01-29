/*
Copyright 2022 The OpenYurt Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package kolecontroller

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
	"github.com/openyurtio/kole/pkg/client/clientset/versioned"
	externalV1alpha1 "github.com/openyurtio/kole/pkg/client/informers/externalversions/lite/v1alpha1"
	listV1alpha1 "github.com/openyurtio/kole/pkg/client/listers/lite/v1alpha1"
)

type QueryNodeController struct {
	kubeclient versioned.Interface
	queue      workqueue.RateLimitingInterface    //workqueue 的引用
	informer   externalV1alpha1.QueryNodeInformer // Informer 的引用
	infEdgeCtl *InfEdgeController
	lister     listV1alpha1.QueryNodeLister
}

// NewQueryNodeController creates a new  QueryNodeController.
func NewQueryNodeController(client versioned.Interface,
	informer externalV1alpha1.QueryNodeInformer,
	infedgeCtl *InfEdgeController) (*QueryNodeController, error) {

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	dsc := &QueryNodeController{
		kubeclient: client,
		informer:   informer,
		queue:      queue,
		lister:     informer.Lister(),
		infEdgeCtl: infedgeCtl,
	}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    dsc.addQueryNode,
		UpdateFunc: dsc.updateQueryNode,
		DeleteFunc: dsc.deleteQueryNode,
	})

	return dsc, nil
}

func (c *QueryNodeController) addQueryNode(obj interface{}) {
	ds := obj.(*v1alpha1.QueryNode)
	klog.V(4).Infof("Adding QueryNode %s", ds.Name)
	c.enqueue(ds)
}

func (c *QueryNodeController) updateQueryNode(oldObj, newObj interface{}) {
	//	oldds := oldObj.(*v1alpha1.QueryNode)
	ds := newObj.(*v1alpha1.QueryNode)
	klog.V(4).Infof("Update QueryNode %s", ds.Name)

	c.enqueue(ds)
}

func (c *QueryNodeController) deleteQueryNode(obj interface{}) {
	ds := obj.(*v1alpha1.QueryNode)
	klog.V(4).Infof("Delete QueryNode %s", ds.Name)
	c.enqueue(ds)
}

func (dsc *QueryNodeController) enqueue(ds *v1alpha1.QueryNode) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(ds)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for object %#v: %v", ds, err))
		return
	}

	// TODO: Handle overlapping controllers better. See comment in ReplicationManager.
	dsc.queue.Add(key)
}

func (c *QueryNodeController) Run(threadiness int, stopCh chan struct{}) {
	defer utilruntime.HandleCrash()

	defer c.queue.ShutDown()

	klog.Info("Starting QueryNode controller")

	// 启动多个 worker 处理 workqueue 中的对象
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	klog.Warningf("Stopping QueryNode controller")
}

func (c *QueryNodeController) runWorker() {
	// 启动无限循环，接收并处理消息
	for c.processNextItem() {

	}
}

// 从 workqueue 中获取对象，并打印信息。
func (c *QueryNodeController) processNextItem() bool {
	key, shutdown := c.queue.Get()
	// 退出
	if shutdown {
		return false
	}

	// 标记此key已经处理
	defer c.queue.Done(key)

	err := c.syncProcess(key.(string))

	c.handleErr(err, key)
	return true
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *QueryNodeController) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing pod %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	utilruntime.HandleError(err)
	klog.Infof("Dropping InfDaemonSet %q out of the queue: %v", key, err)
}

// 获取 key 对应的 object，并打印相关信息
func (c *QueryNodeController) syncProcess(key string) error {

	startTime := time.Now()
	defer func() {
		klog.V(4).Infof("Finished syncing QueryNode set %q (%v)", key, time.Since(startTime))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	qnode, err := c.lister.QueryNodes(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(3).Infof("QueryNode has been deleted %v", key)
		// NEED TO DO

		return nil
	}
	if err != nil {
		return fmt.Errorf("unable to retrieve ds %v from store: %v", key, err)
	}

	s := c.infEdgeCtl.QueryNodeStatusCache.GetNodeStatus(qnode.Spec.NodeName)
	if qnode.Status == nil {
		qnode.Status = make([]*v1alpha1.QueryNodeStatus, 0, 10)
	}

	if !reflect.DeepEqual(s, qnode.Status) {
		qnode.Status = s
		_, err = c.kubeclient.LiteV1alpha1().QueryNodes(namespace).UpdateStatus(context.Background(), qnode, metav1.UpdateOptions{})
		if err != nil {
			klog.Errorf("Update QueryNode error %v", err)
			return err
		}
	}
	return nil
}

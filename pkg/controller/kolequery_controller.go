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
package controller

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

type KoleQueryController struct {
	kubeclient versioned.Interface
	queue      workqueue.RateLimitingInterface    //workqueue 的引用
	informer   externalV1alpha1.KoleQueryInformer // Informer 的引用
	koleCtl    *KoleController
	lister     listV1alpha1.KoleQueryLister
}

// NewKoleQueryController creates a new  KoleQueryController.
func NewKoleQueryController(client versioned.Interface,
	informer externalV1alpha1.KoleQueryInformer,
	koleCtl *KoleController) (*KoleQueryController, error) {

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	kqc := &KoleQueryController{
		kubeclient: client,
		informer:   informer,
		queue:      queue,
		lister:     informer.Lister(),
		koleCtl:    koleCtl,
	}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    kqc.addKoleQuery,
		UpdateFunc: kqc.updateKoleQuery,
		DeleteFunc: kqc.deleteKoleQuery,
	})

	return kqc, nil
}

func (c *KoleQueryController) addKoleQuery(obj interface{}) {
	kq := obj.(*v1alpha1.KoleQuery)
	klog.V(4).Infof("Adding KoleQuery %s", kq.Name)
	c.enqueue(kq)
}

func (c *KoleQueryController) updateKoleQuery(oldObj, newObj interface{}) {
	//	oldds := oldObj.(*v1alpha1.KoleQuery)
	kq := newObj.(*v1alpha1.KoleQuery)
	klog.V(4).Infof("Update KoleQuery %s", kq.Name)
	c.enqueue(kq)
}

func (c *KoleQueryController) deleteKoleQuery(obj interface{}) {
	kq := obj.(*v1alpha1.KoleQuery)
	klog.V(4).Infof("Delete KoleQuery %s", kq.Name)
	c.enqueue(kq)
}

func (c *KoleQueryController) enqueue(ds *v1alpha1.KoleQuery) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(ds)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for object %#v: %v", ds, err))
		return
	}
	c.queue.Add(key)
}

func (c *KoleQueryController) Run(threadiness int, stopCh chan struct{}) {
	defer utilruntime.HandleCrash()

	defer c.queue.ShutDown()

	klog.Info("Starting KoleQuery controller")

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	klog.Warningf("Stopping KoleQuery controller")
}

func (c *KoleQueryController) runWorker() {
	for c.processNextItem() {
	}
}

func (c *KoleQueryController) processNextItem() bool {
	key, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	defer c.queue.Done(key)

	err := c.syncProcess(key.(string))

	c.handleErr(err, key)
	return true
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *KoleQueryController) handleErr(err error, key interface{}) {
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

func (c *KoleQueryController) syncProcess(key string) error {

	startTime := time.Now()
	defer func() {
		klog.V(4).Infof("Finished syncing KoleQuery set %q (%v)", key, time.Since(startTime))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	kq, err := c.lister.KoleQueries(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(3).Infof("KoleQuery has been deleted %v", key)
		// NEED TO DO

		return nil
	}
	if err != nil {
		return fmt.Errorf("unable to retrieve ds %v from store: %v", key, err)
	}

	if kq.Spec.ObjectType == v1alpha1.KoleObjectNode && kq.Spec.ObjectName != "" {
		s := c.koleCtl.QueryNodeStatusCache.GetNodeStatus(kq.Spec.ObjectName)
		if s != nil {
			if kq.Status == nil {
				kq.Status = make([]*v1alpha1.KoleQueryStatus, 0, 0)
			}
			if len(kq.Status) == 0 {
				kq.Status = append(kq.Status, &v1alpha1.KoleQueryStatus{})
			}
			ts := metav1.Now()
			s.LastObservedTime = ts
			kq.Status[0].LastObservedTime = ts
			if !reflect.DeepEqual(s, kq.Status[0]) {
				kq.Status[0] = s
				_, err = c.kubeclient.LiteV1alpha1().KoleQueries(namespace).UpdateStatus(context.Background(), kq, metav1.UpdateOptions{})
				if err != nil {
					klog.Errorf("Update KoleQuery error %v", err)
					return err
				}
			}
		}
	}
	return nil
}

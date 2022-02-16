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
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
	"github.com/openyurtio/kole/pkg/client/clientset/versioned"
	externalV1alpha1 "github.com/openyurtio/kole/pkg/client/informers/externalversions/lite/v1alpha1"
	listV1alpha1 "github.com/openyurtio/kole/pkg/client/listers/lite/v1alpha1"
	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/util"
)

const KoleDaemonSetHashKey = "openyurt.io.koledaemonset/podspec.hash"

type KoleDaemonSetController struct {
	kubeclient versioned.Interface
	queue      workqueue.RateLimitingInterface
	informer   externalV1alpha1.KoleDaemonSetInformer
	koleCtl    *KoleController
	lister     listV1alpha1.KoleDaemonSetLister
}

func Md5PodSpec(obj *v1alpha1.PodSpec) (string, error) {

	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil)), nil
}

// NewKoleDaemonSetController creates a new KoleDaemonSetController.
func NewKoleDaemonSetController(client versioned.Interface, informer externalV1alpha1.KoleDaemonSetInformer, koleCtl *KoleController) (*KoleDaemonSetController, error) {

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	/*
		// This custom indexer will index pods based on their NodeName which will decrease the amount of pods we need to get in simulate() call.
		informer.GetIndexer().AddIndexers(cache.Indexers{
			"nodeName": indexByPodNodeName,
		})
	*/

	dsc := &KoleDaemonSetController{
		kubeclient: client,
		informer:   informer,
		queue:      queue,
		lister:     informer.Lister(),
		koleCtl:    koleCtl,
	}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    dsc.addKoleDaemonSet,
		UpdateFunc: dsc.updateKoleDaemonSet,
		DeleteFunc: dsc.deleteKoleDaemonSet,
	})

	return dsc, nil
}

func (c *KoleDaemonSetController) AddHost(hostName string) {
	klog.V(4).Infof("Adding Host %s", hostName)
	// topic/ pod
	needPublish := make([]*data.Pod, 0, 10240)

	ids, err := c.lister.List(labels.Everything())
	if err != nil {
		return
	}

	c.koleCtl.DesiredPodsCache.SafeOperate(func() {
		if _, ok := c.koleCtl.DesiredPodsCache.Cache[hostName]; !ok {
			desiredPods := make(map[string]*data.Pod)
			for _, ds := range ids {
				// Todo add node selector
				podKey := generateKoleDaemonSetPodKey(ds)
				hash, err := Md5PodSpec(ds.Spec)
				if err != nil {
					klog.Errorf("Generage pod spec hash error %v", err)
					continue
				}
				np := &data.Pod{
					Hash:      hash,
					Name:      generateKoleDaemonSetPodName(ds),
					NameSpace: ds.Namespace,
					Spec:      ds.Spec,
				}
				desiredPods[podKey] = np
				needPublish = append(needPublish, np)
			}
			c.koleCtl.DesiredPodsCache.Cache[hostName] = desiredPods
		}
	})

	for _, ds := range ids {
		c.enqueue(ds)
	}

	go func() {
		topic := filepath.Join(util.TopicDataPrefix, hostName)
		for _, p := range needPublish {
			if err := c.koleCtl.MessageHandler.PublishData(context.Background(), topic, 0, false, p); err != nil {
				klog.Errorf("Mqtt5 publish error %v", err)
				return
			}
		}
	}()
}

func generateKoleDaemonSetPodKey(ds *v1alpha1.KoleDaemonSet) string {
	return fmt.Sprintf("%s-%s", ds.Namespace, generateKoleDaemonSetPodName(ds))
}

func generateKoleDaemonSetPodName(ds *v1alpha1.KoleDaemonSet) string {
	return fmt.Sprintf("koledaemonset-%s", ds.Name)
}

func (c *KoleDaemonSetController) addUpdateKoleDaemonSet(ds *v1alpha1.KoleDaemonSet) {

	needPublish := make(map[string][]*data.Pod)

	podKey := generateKoleDaemonSetPodKey(ds)
	hash, err := Md5PodSpec(ds.Spec)
	if err != nil {
		klog.Errorf("Generage pod spec hash error %v", err)
		return
	}

	c.koleCtl.DesiredPodsCache.WriteRange(func(nodeName string, desiredPodsMap map[string]*data.Pod) {
		// Todo add node selector
		newP := &data.Pod{
			Hash:      hash,
			Name:      generateKoleDaemonSetPodName(ds),
			NameSpace: ds.Namespace,
			Spec:      ds.Spec,
		}
		desiredPodsMap[podKey] = newP
		topic := filepath.Join(util.TopicDataPrefix, nodeName)
		if _, ok := needPublish[topic]; !ok {
			needPublish[topic] = make([]*data.Pod, 0, 10)
		}
		needPublish[topic] = append(needPublish[topic], newP)
	})

	go func() {

		for topic, podList := range needPublish {
			for i, _ := range podList {
				if err := c.koleCtl.MessageHandler.PublishData(context.Background(), topic, 0, false, podList[i]); err != nil {
					klog.Errorf("Mqtt5 publish error %v", err)
					return
				}
			}
		}

	}()
}
func (c *KoleDaemonSetController) addKoleDaemonSet(obj interface{}) {
	ds := obj.(*v1alpha1.KoleDaemonSet)
	klog.Infof("Adding KoleDaemonSet %s time %d", ds.Name, time.Now().Unix())
	c.addUpdateKoleDaemonSet(ds)
	c.enqueue(ds)
}

func (c *KoleDaemonSetController) updateKoleDaemonSet(oldObj, newObj interface{}) {
	oldds := oldObj.(*v1alpha1.KoleDaemonSet)
	ds := newObj.(*v1alpha1.KoleDaemonSet)
	klog.V(4).Infof("Update KoleDaemonSet %s", ds.Name)

	oldHash, err := Md5PodSpec(oldds.Spec)
	if err != nil {
		klog.Errorf("Generage pod spec hash error %v", err)
		return
	}

	newHash, err := Md5PodSpec(ds.Spec)
	if err != nil {
		klog.Errorf("Generage pod spec hash error %v", err)
		return
	}
	if oldHash != newHash {
		c.addUpdateKoleDaemonSet(ds)
	}

	c.enqueue(ds)
}

func (c *KoleDaemonSetController) deleteKoleDaemonSet(obj interface{}) {
	ds := obj.(*v1alpha1.KoleDaemonSet)
	klog.V(4).Infof("Delete KoleDaemonSet %s", ds.Name)

	c.koleCtl.DesiredPodsCache.WriteRange(func(nodeName string, desiredPodsMap map[string]*data.Pod) {
		// Todo add node selector
		podKey := generateKoleDaemonSetPodKey(ds)
		klog.V(4).Infof("Delete KoleDaemonSet pod from node %s , pod key %s", nodeName, podKey)
		oldP := desiredPodsMap[podKey]
		deleteT := metav1.Now()

		deletePod := &data.Pod{
			Hash:            oldP.Hash,
			Name:            oldP.Name,
			NameSpace:       oldP.NameSpace,
			DeleteTimeStamp: &deleteT,
		}

		dataTopic := filepath.Join(util.TopicDataPrefix, nodeName)

		go func() {
			if err := c.koleCtl.MessageHandler.PublishData(context.Background(), dataTopic, 0, false, deletePod); err != nil {
				klog.Errorf("Mqtt5 publish error %v", err)
				return
			}
		}()

		delete(desiredPodsMap, podKey)
	})

	c.enqueue(ds)
}

func (dsc *KoleDaemonSetController) enqueue(ds *v1alpha1.KoleDaemonSet) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(ds)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for object %#v: %v", ds, err))
		return
	}

	// TODO: Handle overlapping controllers better. See comment in ReplicationManager.
	dsc.queue.Add(key)
}

func (c *KoleDaemonSetController) Run(threadiness int, stopCh chan struct{}) {
	defer utilruntime.HandleCrash()

	defer c.queue.ShutDown()

	klog.Info("Starting KoleDaemonSet controller")

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	klog.Warningf("Stopping KoleDaemonSet controller")
}

func (c *KoleDaemonSetController) runWorker() {
	for c.processNextItem() {
	}
}

func (c *KoleDaemonSetController) processNextItem() bool {
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
func (c *KoleDaemonSetController) handleErr(err error, key interface{}) {
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
	klog.Infof("Dropping KoleDaemonSet %q out of the queue: %v", key, err)
}

func (c *KoleDaemonSetController) syncProcess(key string) error {

	startTime := time.Now()
	defer func() {
		klog.V(4).Infof("Finished syncing infdaemon set %q (%v)", key, time.Since(startTime))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	ds, err := c.lister.KoleDaemonSets(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.Warningf("KoleDaemonSet has been deleted %v", key)
		return nil
	}
	if err != nil {
		return fmt.Errorf("unable to retrieve ds %v from store: %v", key, err)
	}

	podKey := generateKoleDaemonSetPodKey(ds)
	hash, err := Md5PodSpec(ds.Spec)
	if err != nil {
		klog.Errorf("Generage pod spec hash error %v", err)
		return err
	}
	var currentNumberScheduled, podready, desirednum int

	desirednum = c.koleCtl.DesiredPodsCache.Len()

	c.koleCtl.ObserverdPodsCache.ReadRange(func(nodeName string, observerdPods map[string]*data.HeartBeatPod) {
		for key, pod := range observerdPods {
			if podKey == key && hash == pod.Hash {
				currentNumberScheduled++
				if pod.Status.Phase == data.HeartBeatPodStatusRunning {
					podready++
				}
			}
		}
	})

	needUpdate := false
	if ds.Status == nil {
		ds.Status = &v1alpha1.KoleDaemonSetStatus{}
		needUpdate = true
	}

	if ds.Status.CurrentNumberScheduled != currentNumberScheduled ||
		ds.Status.NumberReady != podready ||
		ds.Status.DesiredNumberScheduled != desirednum {

		ds.Status.CurrentNumberScheduled = currentNumberScheduled
		ds.Status.NumberReady = podready
		ds.Status.DesiredNumberScheduled = desirednum
		needUpdate = true
	}

	if needUpdate {
		_, err = c.kubeclient.LiteV1alpha1().KoleDaemonSets(ds.Namespace).UpdateStatus(context.Background(), ds, metav1.UpdateOptions{})
		if err != nil {
			klog.Errorf("Update KoleDaemonSet error %v", err)
			return err
		}
	}
	return nil
}

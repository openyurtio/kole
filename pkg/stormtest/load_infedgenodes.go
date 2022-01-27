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
package stormtest

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/client/clientset/versioned"
)

type LoadInfEdgeNodes struct {
	KubeConfig string
	NS         string
	LiteClient versioned.Interface
	PatchNum   int
}

func NewLoadInfEdgeNodes(config *options.LoadInfEdgeNodeFlags) (*LoadInfEdgeNodes, error) {

	c, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		return nil, err
	}
	c.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(2000, 4000)

	crdclient, err := versioned.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	t := &LoadInfEdgeNodes{
		KubeConfig: config.KubeConfig,
		NS:         config.NS,
		LiteClient: crdclient,
		PatchNum:   config.PatchNum,
	}
	return t, nil
}

func (n *LoadInfEdgeNodes) runGood() error {
	klog.V(4).Infof("Start to load all infEdgeNodes ...")
	var timeoutS int64 = 60
	var allNum int
	//var allDataSize int
	var continueStr string
	//	recivedCache := make(map[string]struct{})
	now := time.Now()
	for {
		singleNow := time.Now()
		listedNodes, err := n.LiteClient.LiteV1alpha1().InfEdgeNodes(n.NS).List(context.Background(), metav1.ListOptions{
			TimeoutSeconds: &timeoutS,
			Limit:          int64(n.PatchNum),
			Continue:       continueStr,
		})
		if err != nil {
			klog.Errorf("List all InfEdgeNode in ns[%s] error %v", n.NS, err)
			return err
		}
		getLen := len(listedNodes.Items)
		allNum = allNum + getLen
		singleneedTime := time.Now().Sub(singleNow).Milliseconds()
		klog.V(4).Infof("Get %d infedgenodes need %d ms", getLen, singleneedTime)

		//klog.Infof("Current get list node chunk nums %d ", len(listedNodes.Items))
		continueStr = listedNodes.GetContinue()

		if len(continueStr) == 0 {
			klog.V(4).Infof("continueStr is null , break ")
			break
		}

		if getLen < n.PatchNum {
			break
		}
	}
	needTime := time.Now().Sub(now).Milliseconds()
	klog.Infof("######## List (%d) nodes, patch num %d  need %d ms...",
		allNum,
		n.PatchNum,
		needTime)

	return nil
}
func (n *LoadInfEdgeNodes) Run() error {
	return n.runGood()
}

/*
	// 如果数据量太大， 可能list 同样有问题
	factory := externalversions.NewSharedInformerFactoryWithOptions(n.LiteClient, 0, externalversions.WithNamespace(n.NS))
	now := time.Now()
	klog.Infof("Start infedge nodes informer ...")
	informer := factory.Lite().V1alpha1().InfEdgeNodes().Informer()
	go factory.Start(wait.NeverStop)

	// sync apiserver
	if !cache.WaitForCacheSync(wait.NeverStop,
		informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return fmt.Errorf("time out")
	}
*/

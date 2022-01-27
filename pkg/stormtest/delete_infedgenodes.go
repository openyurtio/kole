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
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/client/clientset/versioned"
)

type DeleteInfEdgeNodes struct {
	KubeConfig     string
	NS             string
	LiteClient     versioned.Interface
	DeletePatchNum int
}

func NewDeleteInfEdgeNodes(config *options.DeleteInfEdgeNodesFlags) (*DeleteInfEdgeNodes, error) {

	c, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		return nil, err
	}
	c.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(20000, 40000)

	crdclient, err := versioned.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	t := &DeleteInfEdgeNodes{
		KubeConfig:     config.KubeConfig,
		NS:             config.NS,
		LiteClient:     crdclient,
		DeletePatchNum: config.DeletePatchNum,
	}
	return t, nil
}

func (n *DeleteInfEdgeNodes) Run() error {

	now := time.Now()
	var timeoutS int64 = 60
	//var allDataSize int
	var continueStr string
	totalNum := 0

	deleteNode := func(name string) {
		err := n.LiteClient.LiteV1alpha1().InfEdgeNodes(n.NS).Delete(context.Background(), name, metav1.DeleteOptions{})
		if err != nil {
			klog.Errorf("Delete InfEdgeNode[%s][%s] error %v", n.NS, name, err)
			return
		}
	}

	for {
		nodes, err := n.LiteClient.LiteV1alpha1().InfEdgeNodes(n.NS).List(context.Background(), metav1.ListOptions{
			TimeoutSeconds: &timeoutS,
			Limit:          int64(n.DeletePatchNum),
			Continue:       continueStr,
		})
		if err != nil {
			klog.Errorf("List InfEdgeNode from ns %s error %v", n.NS, err)
			return err
		}
		continueStr = nodes.GetContinue()
		group := sync.WaitGroup{}
		for _, s := range nodes.Items {
			totalNum++
			group.Add(1)
			go func(name string) {
				defer group.Done()
				deleteNode(name)
			}(s.Name)
		}
		group.Wait()

		if len(nodes.Items) < n.DeletePatchNum {
			break
		}

	}
	needTime := time.Now().Sub(now).Milliseconds()
	klog.Infof("######## Delete %d InfEdgeNode (patch %d) success, need %d ms ", totalNum, n.DeletePatchNum, needTime)
	return nil
}

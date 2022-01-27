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
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
	"github.com/openyurtio/kole/pkg/client/clientset/versioned"
)

type CreateInfEdgeNodes struct {
	KubeConfig  string
	NodeNum     int
	BatchNum    int
	IsSmallSize bool
	NS          string
	LiteClient  versioned.Interface
}

func NewCreateInfEdgeNodes(config *options.CreateInfEdgeNodesFlags) (*CreateInfEdgeNodes, error) {

	c, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		return nil, err
	}
	c.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(20000, 40000)

	crdclient, err := versioned.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	t := &CreateInfEdgeNodes{
		KubeConfig:  config.KubeConfig,
		NodeNum:     config.NodeNum,
		BatchNum:    config.BatchNum,
		NS:          config.NS,
		LiteClient:  crdclient,
		IsSmallSize: config.IsSmallSize,
	}
	return t, nil
}

func (n *CreateInfEdgeNodes) Run() error {
	//var allSize uint64
	now := time.Now()
	var tmpNode *v1alpha1.InfEdgeNode
	if n.IsSmallSize {
		tmpNode = smallInfEdgeNodeTemplate
		tmpNode.Namespace = n.NS
		//singleSize = SmallInfEdgeNodeTemplateSize
	} else {
		tmpNode = infEdgeNodeTemplate
		tmpNode.Namespace = n.NS
		//singleSize = InfEdgeNodeTemplateSize
	}

	createNode := func(i, j int, t v1alpha1.InfEdgeNode) {
		flag := fmt.Sprintf("%s", uuid.New())
		tt := &t
		tt.Name = fmt.Sprintf("%d-%d-%s", i, j, flag)
		tt.ResourceVersion = ""
		for i := 0; i < 5; i++ {
			_, err := n.LiteClient.LiteV1alpha1().InfEdgeNodes(tt.GetNamespace()).Create(context.Background(), tt, metav1.CreateOptions{})
			if err != nil {
				klog.Errorf("create infEdgeNode [%s][%s] error %v, prepare to continue %d...", tt.GetNamespace(), tt.GetName(), err, i)
				time.Sleep(time.Second * 5)
				continue
			} else {
				break
			}
		}
	}

	if n.NodeNum/n.BatchNum == 0 {
		group := sync.WaitGroup{}
		for i := 0; i < n.NodeNum; i++ {
			group.Add(1)
			go func(i, j int, t v1alpha1.InfEdgeNode) {
				defer group.Done()
				createNode(i, j, t)
			}(i, 0, *tmpNode)
		}
		group.Wait()
	} else {
		i := 0
		for ; i < n.NodeNum/n.BatchNum; i++ {
			group := sync.WaitGroup{}
			for j := 0; j < n.BatchNum; j++ {
				group.Add(1)
				go func(i, j int, t v1alpha1.InfEdgeNode) {
					defer group.Done()
					createNode(i, j, t)
				}(i, j, *tmpNode)
			}
			group.Wait()
		}

		group := sync.WaitGroup{}
		for j := i; j < i+n.NodeNum%n.BatchNum; j++ {
			group.Add(1)
			go func(i, j int, t v1alpha1.InfEdgeNode) {
				defer group.Done()
				createNode(i, j, t)
			}(j, 0, *tmpNode)
		}
		group.Wait()
	}
	//allSize = allSize / 1024
	needTime := time.Now().Sub(now).Milliseconds()

	klog.V(4).Infof("######## Total create %d InfEdgeNode . need %d ms",
		n.NodeNum,
		needTime,
	)

	return nil
}

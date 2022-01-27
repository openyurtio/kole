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

type LoadQueryNode struct {
	KubeConfig string
	NS         string
	Name       string
	LiteClient versioned.Interface
}

func NewLoadQueryNode(config *options.LoadQueryNodeFlags) (*LoadQueryNode, error) {

	c, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		return nil, err
	}
	c.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(2000, 4000)

	crdclient, err := versioned.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	t := &LoadQueryNode{
		KubeConfig: config.KubeConfig,
		NS:         config.Ns,
		LiteClient: crdclient,
		Name:       config.Name,
	}
	return t, nil
}

func (n *LoadQueryNode) Run() error {
	klog.V(4).Infof("Start to load all cr %s ...", n.Name)
	var num int
	now := time.Now()
	for {
		singleNow := time.Now()
		_, err := n.LiteClient.LiteV1alpha1().QueryNodes(n.NS).Get(context.Background(), n.Name, metav1.GetOptions{})
		if err != nil {
			klog.Errorf("Get querynode in ns[%s] error %v", n.NS, err)
			return err
		}

		singleneedTime := time.Now().Sub(singleNow).Milliseconds()
		klog.Infof("[%d]Get querynode cr %s in %s   need %d ms", num, n.Name, n.NS, singleneedTime)
		num++

		if num < 1000 {
			break
		}
	}
	needTime := time.Now().Sub(now).Milliseconds()
	klog.Infof("######## Get 1000 querynode,  average time is %d ms...",
		needTime/1000)

	return nil
}

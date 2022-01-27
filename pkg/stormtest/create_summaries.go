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

type CreateSummaries struct {
	KubeConfig string
	SumNum     int
	BatchNum   int
	NS         string
	LiteClient versioned.Interface
}

func NewCreateSummaries(config *options.CreateSummariesFlags) (*CreateSummaries, error) {

	c, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		return nil, err
	}
	// set rate limit
	c.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(10000, 10000)

	crdclient, err := versioned.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	t := &CreateSummaries{
		KubeConfig: config.KubeConfig,
		SumNum:     config.SumNum,
		BatchNum:   config.BatchNum,
		NS:         config.NS,
		LiteClient: crdclient,
	}
	return t, nil
}

var tmpSummarisData []byte

func init() {
	tmpSummarisData = make([]byte, 1000*1000, 1000*1000)
}
func (n *CreateSummaries) Run() error {
	//var allSize uint64
	now := time.Now()
	tmpSum := &v1alpha1.Summary{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: n.NS,
			Name:      "",
		},
		Data: tmpSummarisData,
	}

	createSummary := func(i, j int, t v1alpha1.Summary) {
		flag := fmt.Sprintf("%s", uuid.New())
		tt := &t
		tt.Name = fmt.Sprintf("%d-%d-%s", i, j, flag)
		tt.ResourceVersion = ""
		klog.V(4).Infof("prepare to create summaris [%s][%s] ...", tt.Namespace, tt.Name)
		_, err := n.LiteClient.LiteV1alpha1().Summaries(tt.GetNamespace()).Create(context.Background(), tt, metav1.CreateOptions{})
		if err != nil {
			klog.Errorf("create summaris [%s][%s] error %v ", tt.GetNamespace(), tt.GetName(), err)
			return
		}
		klog.V(4).Infof("create summaris [%s][%s] success", tt.Namespace, tt.Name)
	}
	klog.Infof("Begin to create summaris ...")

	if n.SumNum/n.BatchNum == 0 {
		group := sync.WaitGroup{}
		for i := 0; i < n.SumNum; i++ {
			group.Add(1)
			go func(i, j int, t v1alpha1.Summary) {
				defer group.Done()
				createSummary(i, j, t)
			}(i, 0, *tmpSum)
		}
		group.Wait()
	} else {
		i := 0
		for ; i < n.SumNum/n.BatchNum; i++ {
			group := sync.WaitGroup{}
			for j := 0; j < n.BatchNum; j++ {
				group.Add(1)
				go func(i, j int, t v1alpha1.Summary) {
					defer group.Done()
					createSummary(i, j, t)
				}(i, j, *tmpSum)
			}
			group.Wait()
		}

		group := sync.WaitGroup{}
		for j := i; j < i+n.SumNum%n.BatchNum; j++ {
			group.Add(1)
			go func(i, j int, t v1alpha1.Summary) {
				defer group.Done()
				createSummary(i, j, t)
			}(j, 0, *tmpSum)
		}
		group.Wait()
	}
	//allSize = allSize / 1024
	needTime := time.Now().Sub(now).Milliseconds()

	klog.Infof("######## Total create %d Summaris. need %d ms",
		n.SumNum,
		needTime,
	)
	return nil
}

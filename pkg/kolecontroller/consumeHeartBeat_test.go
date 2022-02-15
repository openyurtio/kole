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
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"

	"github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
	"github.com/openyurtio/kole/pkg/client/clientset/versioned"
	"github.com/openyurtio/kole/pkg/client/informers/externalversions"
	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/util"
)

func TestConsumeHeartBeat(t *testing.T) {
	stop := make(chan struct{}, 1)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Get HomeDir error %v", err)
	}

	c, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir, ".kube/config"))
	if err != nil {
		t.Fatalf("Build config error %v", err)
	}

	// set rate limit
	c.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(10000, 10000)

	crdclient, err := versioned.NewForConfig(c)
	if err != nil {
		t.Fatalf("Build versioned config error %v", err)
	}

	koleInstance := &KoleController{
		HeartBeatCache: &HeartBeatCache{
			RWMutex: &sync.RWMutex{},
			Cache:   make(map[string]*data.HeartBeat),
		},
		HeartBeatFilter: &HeartBeatFilter{
			Mutex:  &sync.Mutex{},
			Filter: make(map[string]*FilterInfo),
		},

		ObserverdPodsCache: &ObserverdPodsCache{
			RWMutex: &sync.RWMutex{},
			Cache:   make(map[string]map[string]*data.HeartBeatPod),
		},
		QueryNodeStatusCache: &QueryNodeStatusCache{
			RWMutex:      &sync.RWMutex{},
			NameToStatus: make(map[string]*v1alpha1.QueryNodeStatus),
		},
		DesiredPodsCache: &DesiredPodsCache{
			RWMutex: &sync.RWMutex{},
			Cache:   make(map[string]map[string]*data.Pod)},
	}
	factory := externalversions.NewSharedInformerFactory(crdclient, time.Second*70)
	infDaemonSetInfor := factory.Lite().V1alpha1().InfDaemonSets()
	controller, err := NewInfDaemonSetController(crdclient, infDaemonSetInfor, infedge)

	go factory.Start(stop)

	if !cache.WaitForCacheSync(wait.NeverStop,
		infDaemonSetInfor.Informer().HasSynced,
	) {
		t.Fatalf("Wait for cache sync error %v", err)
	}

	go controller.Run(5, stop)

	koleInstance.InfDaemonSetController = controller

	hb := util.InitMockHeartBeat("", 1)

	for i := 1; i <= 10; i++ {
		n := time.Now()
		for j := 0; j < 100000; j++ {
			hb.Name = fmt.Sprintf("%d-%d", i, j)
			koleInstance.ConsumeSingleHeartBeat(hb)
		}
		t.Logf("%d0w use %d ms", i, time.Now().Sub(n).Milliseconds())
	}
}

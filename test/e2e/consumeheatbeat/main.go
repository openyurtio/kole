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
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"

	"github.com/openyurtio/kole/pkg/client/clientset/versioned"
	"github.com/openyurtio/kole/pkg/client/informers/externalversions"
	"github.com/openyurtio/kole/pkg/controller"
	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/util"
)

var level int

func init() {
	flag.IntVar(&level, "level", level, "number level")
}

func main() {
	flag.Parse()

	stop := make(chan struct{}, 1)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Get HomeDir error %v", err)
	}

	c, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir, ".kube/config"))
	if err != nil {
		log.Fatalf("Build config error %v", err)
	}

	// set rate limit
	c.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(10000, 10000)

	crdclient, err := versioned.NewForConfig(c)
	if err != nil {
		log.Fatalf("Build versioned config error %v", err)
	}

	infedge := &controller.InfEdgeController{

		HeatBeatCache: &controller.HeatBeatCache{
			RWMutex: &sync.RWMutex{},
			Cache:   make(map[string]*data.HeatBeat),
		},
		HeatBeatFilter: &controller.HeatBeatFilter{
			Mutex:  &sync.Mutex{},
			Filter: make(map[string]*controller.FilterInfo),
		},
		ObserverdPodsCache: &controller.ObserverdPodsCache{
			RWMutex: &sync.RWMutex{},
			Cache:   make(map[string]map[string]*data.HeatBeatPod),
		},
		DesiredPodsCache: &controller.DesiredPodsCache{
			RWMutex: &sync.RWMutex{},
			Cache:   make(map[string]map[string]*data.Pod)},
	}
	factory := externalversions.NewSharedInformerFactory(crdclient, time.Second*70)
	infDaemonSetInfor := factory.Lite().V1alpha1().InfDaemonSets()
	controller, err := controller.NewInfDaemonSetController(crdclient, infDaemonSetInfor, infedge)

	go factory.Start(stop)

	if !cache.WaitForCacheSync(wait.NeverStop,
		infDaemonSetInfor.Informer().HasSynced,
	) {
		log.Fatalf("Wait for cache sync error %v", err)
	}

	go controller.Run(5, stop)

	infedge.InfDaemonSetController = controller

	hb := util.InitMockHeatBeat("", 1)

	pods := make([]*data.HeatBeatPod, 0, 10)

	i := 1
	for ; i <= 5; i++ {
		p := &data.HeatBeatPod{
			Hash:      "",
			Name:      fmt.Sprintf("d%d", i),
			NameSpace: "infedge",
			Status: &data.HeatBeatPodStatus{
				Phase: data.HeatBeatPodStatusRunning,
			},
		}
		pods = append(pods, p)
	}

	for ; i <= 10; i++ {
		p := &data.HeatBeatPod{
			Hash:      "",
			Name:      fmt.Sprintf("delete%d", i),
			NameSpace: "infedge",
			Status: &data.HeatBeatPodStatus{
				Phase: data.HeatBeatPodStatusRunning,
			},
		}
		pods = append(pods, p)
	}
	hb.Pods = pods
	fmt.Printf("level %d\n", level)

	for j := 0; j < level; j++ {
		hb.Name = fmt.Sprintf("first-%d", j)
		infedge.ConsumeSingleHeatBeat(hb)
	}

	n := time.Now()
	for i := 0; i < level; i++ {
		hb.Name = fmt.Sprintf("hb-%d", i)
		infedge.ConsumeSingleHeatBeat(hb)
	}
	fmt.Printf("%d add %d use %d ms\n", level, level, time.Now().Sub(n).Milliseconds())
}

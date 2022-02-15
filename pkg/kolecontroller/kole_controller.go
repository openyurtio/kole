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
	"sync"
	"time"

	outmqtt "github.com/eclipse/paho.mqtt.golang"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/kole-controller/app/options"
	"github.com/openyurtio/kole/pkg/client/clientset/versioned"
	"github.com/openyurtio/kole/pkg/client/informers/externalversions"
	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/message"
	"github.com/openyurtio/kole/pkg/util"
)

// +kubebuilder:rbac:groups=lite.openyurt.io,resources=infdaemonsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lite.openyurt.io,resources=infdaemonsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lite.openyurt.io,resources=querynodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lite.openyurt.io,resources=querynodes/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=lite.openyurt.io,resources=summaries,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lite.openyurt.io,resources=summaries/status,verbs=get;update;patch

type SendMessage struct {
	Topic string
	Data  interface{}
}

type KoleController struct {
	MessageHandler message.MessageHandler

	ObserverdPodsCache   *ObserverdPodsCache
	DesiredPodsCache     *DesiredPodsCache
	QueryNodeStatusCache *QueryNodeStatusCache

	InfDaemonSetController *InfDaemonSetController
	QueryNodeController    *QueryNodeController

	// key nodename
	HeartBeatCache *HeartBeatCache

	HeartBeatTimeOut int64

	HeartBeatFilter *HeartBeatFilter

	DataProcess       DataProcesser
	SnapshotInterval  int
	SummaryNS         string
	SnapdSummaryNames []string
	LiteClient        versioned.Interface
	LasterSnapIndex   int64
	LasterSnapTime    int64
	FirstSnapTime     int64
	ReceiveNum        int64
}

func NewMainKoleController(stop chan struct{}, config *options.KoleControllerFlags, processer DataProcesser) (*KoleController, error) {

	c, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		return nil, err
	}

	// set rate limit
	c.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(2000, 3000)

	crdclient, err := versioned.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	heatBeatCache, heatBeatFilter, snapedName, observerdPods, nodeStatus, err := LoadSnapShot(crdclient, config, processer)
	if err != nil {
		return nil, err
	}

	koleInstance := &KoleController{
		SummaryNS:         config.NameSpace,
		HeartBeatTimeOut:  int64(config.HBTimeOut),
		LiteClient:        crdclient,
		DataProcess:       processer,
		SnapshotInterval:  config.SnapshotInterval,
		SnapdSummaryNames: snapedName,

		HeartBeatCache: &HeartBeatCache{
			RWMutex: &sync.RWMutex{},
			Cache:   heatBeatCache,
		},
		HeartBeatFilter: &HeartBeatFilter{
			Mutex:  &sync.Mutex{},
			Filter: heatBeatFilter,
		},

		ObserverdPodsCache: &ObserverdPodsCache{
			RWMutex: &sync.RWMutex{},
			Cache:   observerdPods,
		},
		QueryNodeStatusCache: &QueryNodeStatusCache{
			RWMutex:      &sync.RWMutex{},
			NameToStatus: nodeStatus,
		},
		DesiredPodsCache: &DesiredPodsCache{
			RWMutex: &sync.RWMutex{},
			Cache:   make(map[string]map[string]*data.Pod)},
	}

	factory := externalversions.NewSharedInformerFactory(crdclient, time.Second*70)
	infDaemonSetInfor := factory.Lite().V1alpha1().InfDaemonSets()
	controller, err := NewInfDaemonSetController(crdclient, infDaemonSetInfor, koleInstance)

	queryNodeInfor := factory.Lite().V1alpha1().QueryNodes()
	queryNodeController, err := NewQueryNodeController(crdclient, queryNodeInfor, koleInstance)

	go factory.Start(stop)

	if !cache.WaitForCacheSync(wait.NeverStop,
		infDaemonSetInfor.Informer().HasSynced,
		queryNodeInfor.Informer().HasSynced,
	) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return nil, fmt.Errorf("time out")
	}

	go controller.Run(5, stop)
	go queryNodeController.Run(5, stop)

	koleInstance.InfDaemonSetController = controller
	koleInstance.QueryNodeController = queryNodeController

	for nodeName, _ := range heatBeatCache {
		controller.AddHost(nodeName)
	}
	if !config.IsMqtt5 {
		h, err := message.NewMqtt3Handler(config.Mqtt3Flags.MqttBroker, config.Mqtt3Flags.MqttBrokerPort, config.Mqtt3Flags.MqttInstance, config.Mqtt3Flags.MqttGroup,
			"kole-controller",
			map[string]outmqtt.MessageHandler{
				util.TopicHeartBeat: koleInstance.Mqtt3SubEdgeHeartBeat,
			})
		if err != nil {
			return nil, err
		}
		koleInstance.MessageHandler = h
	} else {
		// mqtt 5
		h, err := message.NewMqtt5Handler(config.Mqtt5Flags.MqttServer, koleInstance.Mqtt5CreateSubscribes(), "kolecontroller-mqtt-v5", false)
		if err != nil {
			return nil, err
		}
		koleInstance.MessageHandler = h
	}

	klog.V(4).Infof("Create kole cloud mqtt client successfully")

	return koleInstance, nil
}

func (l *KoleController) Run() error {

	go l.SnapShotLoop()
	return nil
}

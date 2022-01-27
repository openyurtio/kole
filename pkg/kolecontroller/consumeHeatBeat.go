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
	"context"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/util"
)

func (c *InfEdgeController) ConsumeHeatBeatDirect(hb *data.HeatBeat) {

	sync_pods := c.ConsumeSingleHeatBeat(hb)

	dataTopic := filepath.Join(util.TopicDataPrefix, hb.Name)

	go func(topic string, sends []*data.Pod) {
		for i, _ := range sends {
			if err := c.MessageHandler.PublishData(context.Background(), dataTopic, 0, false, sends[i]); err != nil {
				klog.Errorf("Mqtt5 publish error %v", err)
				return
			}
		}
	}(dataTopic, sync_pods)
}

func (c *InfEdgeController) ConsumeSingleHeatBeat(hb *data.HeatBeat) []*data.Pod {
	c.ReceiveNum++

	klog.V(5).Infof("Receive heatbeat Indentifier[%s] Name[%s] State[%s]", hb.Identifier, hb.Name, hb.State)

	if !c.HeatBeatFilter.SetHeatBeat(hb) {
		return []*data.Pod{}
	}

	c.ObserverdPodsCache.SafeSetHeatBeat(hb)

	// 获得的与期望的比较
	sync_pods := make([]*data.Pod, 0, 20)

	c.DesiredPodsCache.SafeOperate(func() {
		desiredPods, ok := c.DesiredPodsCache.Cache[hb.Name]
		if !ok {
			return
		}
		for _, desiredPod := range desiredPods {
			find := false
			needUpdate := false
			for _, hbPod := range hb.Pods {
				if desiredPod.Key() == hbPod.Key() {
					find = true
					if desiredPod.Hash != hbPod.Hash {
						needUpdate = true
					}
					break
				}
			}
			if !find || needUpdate {
				sync_pods = append(sync_pods, &data.Pod{
					Hash:      desiredPod.Hash,
					Name:      desiredPod.Name,
					NameSpace: desiredPod.NameSpace,
					Spec:      desiredPod.Spec,
				})
			}
		}
		deleteT := metav1.Now()
		for _, hbPod := range hb.Pods {
			if _, ok := desiredPods[hbPod.Key()]; !ok {
				sync_pods = append(sync_pods, &data.Pod{
					Hash:            hbPod.Hash,
					Name:            hbPod.Name,
					NameSpace:       hbPod.NameSpace,
					DeleteTimeStamp: &deleteT,
				})
			}
		}
	})

	c.HeatBeatCache.ReceiveHeatBeat(hb, c.InfDaemonSetController)

	return sync_pods
}

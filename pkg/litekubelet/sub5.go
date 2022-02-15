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
package litekubelet

import (
	"path/filepath"

	"github.com/eclipse/paho.golang/paho"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/pkg/cache"
	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/message"
	"github.com/openyurtio/kole/pkg/util"
)

func (c *LiteKubelet) CreateSubscribes5() []*message.SingleSubcribe {

	subs := make([]*message.SingleSubcribe, 0, 2)
	//CTL
	subs = append(subs, &message.SingleSubcribe{
		Topic: filepath.Join(util.TopicCTLPrefix, c.HostnameOverride),
		Option: paho.SubscribeOptions{
			QoS: 1,
		},
		Handler: func(publish *paho.Publish) {
			c.Sub5CtlChan <- publish
		},
	})
	//DATA
	subs = append(subs, &message.SingleSubcribe{
		Topic: filepath.Join(util.TopicDataPrefix, c.HostnameOverride),
		Option: paho.SubscribeOptions{
			QoS: 1,
		},
		Handler: func(publish *paho.Publish) {
			c.Sub5DataChan <- publish
		},
	})

	return subs
}

func (c *LiteKubelet) ConsumeSubLoop() {
	// CTL
	go func() {
		for p := range c.Sub5CtlChan {
			ack, err := data.UnmarshalPayloadToHeartBeatACK(p.Payload)
			if err != nil {
				klog.Errorf("Unmarshalpayload to headbeatack error %v", err)
				return
			}
			cache.GetDefaultTimeoutCache().Set(ack.Identifier, ack)
			klog.V(5).Infof("Sub heatbeat topic %s", p.Topic)
		}
	}()

	//Data
	go func() {
		for p := range c.Sub5DataChan {
			pod, err := data.UnmarshalPayloadToPod(p.Payload)
			if err != nil {
				klog.Errorf("Unmarshalpayload to Pod error %v", err)
				return
			}
			c.ReceivePodDataNum++
			SyncLocalPod(pod)
			klog.V(4).Infof("Sub data topic %s ,data %s", p.Topic, p.Payload)
		}
	}()
}

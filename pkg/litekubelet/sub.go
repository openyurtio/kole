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
	"time"

	outmqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/pkg/cache"
	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/util"
)

func (c *LiteKubelet) RegisterAndSubscribeTopic() error {
	c.SubTopics[filepath.Join(util.TopicCTLPrefix, c.HostnameOverride)] = c.SubCTL
	c.SubTopics[filepath.Join(util.TopicDataPrefix, c.HostnameOverride)] = c.SubData

	for t, f := range c.SubTopics {
		token := c.MqttClient.Subscribe(t, 1, f)
		if token.WaitTimeout(time.Second * 5) {
			if err := token.Error(); err != nil {
				klog.Errorf("Client %s subscribe topic %s error %v", c.HostnameOverride, t, err)
				continue
			} else {
				klog.V(5).Infof("Client %s subscribe topic %s successfully", c.HostnameOverride, t)
			}
		} else {
			klog.Errorf("Client %s subscribe topic %s timeout", c.HostnameOverride, t)
			continue
		}
	}
	return nil
}

func (c *LiteKubelet) SubCTL(client outmqtt.Client, message outmqtt.Message) {
	ack, err := data.UnmarshalPayloadToHeartBeatACK(message.Payload())
	if err != nil {
		klog.Errorf("Unmarshalpayload to headbeatack error %v", err)
		return
	}
	cache.GetDefaultTimeoutCache().Set(ack.Identifier, ack)
	klog.V(5).Infof("Sub heatbeat topic %s", message.Topic())
}

func (c *LiteKubelet) SubData(client outmqtt.Client, message outmqtt.Message) {
	pod, err := data.UnmarshalPayloadToPod(message.Payload())
	if err != nil {
		klog.Errorf("Unmarshalpayload to Pod error %v", err)
		return
	}
	SyncLocalPod(pod)
	klog.V(5).Infof("sub data topic %s", message.Topic())
}

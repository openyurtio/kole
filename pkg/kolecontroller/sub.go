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
	outmqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/pkg/data"
)

func (c *InfEdgeController) Mqtt3SubEdgeHeatBeat(client outmqtt.Client, message outmqtt.Message) {
	//go func(message outmqtt.Message) {
	hb, err := data.UnmarshalPayloadToHeatBeat(message.Payload())
	if err != nil {
		klog.Errorf("UnmarshalPayloadToHeatBeat error %v", err)
		return
	}
	c.ConsumeHeatBeatDirect(hb)
	//c.HeatBeatQueue <- hb
	klog.V(5).Infof("sub heatbeat topic %s Name %s State %s", message.Topic(), hb.Name, hb.State)
	//}(mes)
}

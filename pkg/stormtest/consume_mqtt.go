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

	outmqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/message"
)

type ConsumeMqtt struct {
	MessageNum         int
	ReceivedMessageNum int
	StartTime          time.Time
	StopCh             chan struct{}
	Topic              string
	MessageHandler     message.MessageHandler
}

func NewConsumeMqtt(config *options.ConsumeMqttFlags) (*ConsumeMqtt, error) {

	t := &ConsumeMqtt{
		MessageNum: config.MessageNum,
		StopCh:     make(chan struct{}, 1),
		Topic:      config.Topic,
	}

	h, err := message.NewMqtt3Handler(config.MqttBroker, config.MqttBrokerPort, config.MqttInstance, config.MqttGroup,
		"consume-sub", map[string]outmqtt.MessageHandler{
			t.Topic: t.testSub,
		})
	if err != nil {
		return nil, err
	}
	t.MessageHandler = h

	return t, nil
}

func (n *ConsumeMqtt) testSub(client outmqtt.Client, message outmqtt.Message) {
	hb, err := data.UnmarshalPayloadToHeatBeat(message.Payload())
	if err != nil {
		klog.Errorf("UnmarshalPayloadToHeatBeat error %v", err)
		return
	}
	klog.Infof("Receive Message %s", hb.Name)
	n.ReceivedMessageNum++
	if n.ReceivedMessageNum%100 == 0 {
		klog.Infof("Has receive %d num memssage", n.ReceivedMessageNum)
	}
	//klog.Infof("Receive message ...")
	if n.MessageNum != 0 && n.ReceivedMessageNum >= n.MessageNum {
		n.StopCh <- struct{}{}
	}
}

func (n *ConsumeMqtt) Run(ctx context.Context) error {

	n.StartTime = time.Now()
	for {
		select {
		case <-n.StopCh:
			klog.Infof("Receive %d message need %d ms", n.MessageNum, time.Now().Sub(n.StartTime).Milliseconds())
			return nil
		case <-ctx.Done():
			klog.Infof("Receive %d message ,timeout need %d ms", n.MessageNum, time.Now().Sub(n.StartTime).Milliseconds())
			return nil
		}
	}
}

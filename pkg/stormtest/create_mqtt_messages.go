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
	"time"

	"github.com/openyurtio/kole/pkg/message"

	outmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/data"
)

type CreateMqttMessages struct {
	MessageNum     int
	Topic          string
	MessageHandler message.MessageHandler
}

func NewCreateMqttMessage(config *options.CreateMqttMessageFlags) (*CreateMqttMessages, error) {
	var hostnameOverride string
	if !config.RandClientID {
		hostnameOverride = "stormtest-client-send"
	} else {
		hostnameOverride = fmt.Sprintf("create-rand-%s", uuid.New())
	}

	if len(config.AccessKey) == 0 || len(config.AccessSecret) == 0 {
		klog.Errorf("ACCESS_KEY[%s] or ACCESS_SECRET[%s] is nil", config.AccessKey, config.AccessSecret)
		return nil, fmt.Errorf("access key or secret is nil")
	}

	//c := mqtt.NewSessionMqttClient(config.MqttBroker, config.MqttBrokerPort, clientID, username, passwd)
	t := &CreateMqttMessages{
		MessageNum: config.MessageNum,
		Topic:      config.Topic,
	}

	h, err := message.NewMqtt3Handler(config.MqttBroker, config.MqttBrokerPort, config.MqttInstance, config.MqttGroup, hostnameOverride, map[string]outmqtt.MessageHandler{})
	if err != nil {
		return nil, fmt.Errorf("NewMqtt3Handler error")
	}
	t.MessageHandler = h

	klog.V(4).Infof("create mqtt client successful, id %s", hostnameOverride)

	return t, nil
}

func initHeatBeat() *data.HeatBeat {
	labels := make(map[string]string)
	labels["HostName"] = "testtttttttttttttttttt"

	beat := &data.HeatBeat{
		SeqNum:     0,
		TimeStamp:  time.Now().Unix(),
		Identifier: fmt.Sprintf("%v", uuid.New()),
		State:      data.HeatBeatRegistering,
		Name:       "testtttttttttttttttttt",
		Labels:     labels,
	}

	return beat
}

func (n *CreateMqttMessages) Run() error {
	hb := initHeatBeat()

	send := func(i int) error {
		hb.Name = fmt.Sprintf("name-%d", i)

		if err := n.MessageHandler.PublishData(context.Background(), n.Topic, 0, false, hb); err != nil {
			return err
		}

		return nil
	}
	nn := time.Now()
	for i := 0; i < n.MessageNum; i++ {
		send(i)
		if i%1000 == 0 {
			klog.Infof("Has send %d message ...", i)
		}
	}
	klog.Infof("Send %d message success , need %d ms ...", n.MessageNum, time.Now().Sub(nn).Milliseconds())

	return nil
}

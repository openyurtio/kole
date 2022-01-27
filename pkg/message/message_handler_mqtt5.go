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
package message

import (
	"context"
	"fmt"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"k8s.io/klog/v2"
)

type Mqtt5Handler struct {
	MqttClient     *autopaho.ConnectionManager
	MqttDataClient *autopaho.ConnectionManager
	MqttAckClient  *autopaho.ConnectionManager
	OneClient      bool
}

func NewMqtt5Handler(server string,
	subs []*SingleSubcribe,
	hostnameOverride string,
	oneClient bool) (*Mqtt5Handler, error) {
	var err error
	h := &Mqtt5Handler{
		OneClient: oneClient,
	}
	// mqtt 5
	h.MqttClient, err = NewMqtt5Manager(context.Background(),
		30,
		3600,
		true,
		65535,
		time.Second*5,
		time.Minute*60,
		server, fmt.Sprintf("%s-sub", hostnameOverride), subs)
	if err != nil {
		klog.Errorf("New mqtt sub client error %v", err)
		return nil, err
	}

	if oneClient {
		return h, nil
	}

	h.MqttAckClient, err = NewMqtt5Manager(context.Background(),
		30,
		3600,
		true,
		10000,
		time.Second*5,
		time.Second*1200,
		server, fmt.Sprintf("%s-ack", hostnameOverride), nil)
	if err != nil {
		klog.Errorf("New mqtt pub client error %v", err)
		return nil, err
	}

	h.MqttDataClient, err = NewMqtt5Manager(context.Background(),
		30,
		3600,
		true,
		10000,
		time.Second*5,
		time.Second*1200,
		server, fmt.Sprintf("%s-data", hostnameOverride), nil)
	if err != nil {
		klog.Errorf("New mqtt pub client error %v", err)
		return nil, err
	}

	return h, nil
}
func (m *Mqtt5Handler) PublishAck(ctx context.Context, topic string, qos byte, retained bool, object interface{}) error {
	if m.OneClient {
		return Mqtt5Send(ctx, m.MqttClient, topic, qos, retained, object)
	}
	return Mqtt5Send(ctx, m.MqttAckClient, topic, qos, retained, object)
}

func (m *Mqtt5Handler) PublishData(ctx context.Context, topic string, qos byte, retained bool, object interface{}) error {
	if m.OneClient {
		return Mqtt5Send(ctx, m.MqttClient, topic, qos, retained, object)
	}
	return Mqtt5Send(ctx, m.MqttDataClient, topic, qos, retained, object)
}

var _ MessageHandler = &Mqtt5Handler{}

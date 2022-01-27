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
	"os"
	"time"

	outmqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/pkg/util"
)

type Mqtt3Handler struct {
	MqttSubClient    outmqtt.Client
	MqttDataClient   outmqtt.Client
	MqttAckClient    outmqtt.Client
	SubTopicHandlers map[string]outmqtt.MessageHandler
}

func (l *Mqtt3Handler) Reconnect(client outmqtt.Client) {
	time.Sleep(time.Second * 1)
	l.subscribeTopics()
	klog.V(4).Infof("Connected mqtt broker , resub topic successfull...")
}

func (l *Mqtt3Handler) ReconnectNoSub(client outmqtt.Client) {
	time.Sleep(time.Second * 1)
	reader := client.OptionsReader()
	klog.V(4).Infof("No Subscribe client %s connected mqtt broker successfull ...", reader.ClientID())
}

func (l *Mqtt3Handler) LostConnectHandler(client outmqtt.Client, err error) {
	klog.V(5).Infof("Connect lost:%v", err)
}

func (l *Mqtt3Handler) subscribeTopics() {

	for t, f := range l.SubTopicHandlers {
		token := l.MqttSubClient.Subscribe(t, 1, f)
		token.Wait()
		if err := token.Error(); err != nil {
			klog.Fatalf("Subscribe topic %s error %v", t, err)
		}
		klog.V(4).Infof("Subscribe topic %s successfully", t)
	}
}

func NewMqtt3Handler(
	broker string,
	port int,
	instance,
	group string,
	hostname string,
	subTopicsHandlers map[string]outmqtt.MessageHandler) (*Mqtt3Handler, error) {

	key := os.Getenv("ACCESS_KEY")
	secret := os.Getenv("ACCESS_SECRET")

	if len(key) == 0 || len(secret) == 0 {
		klog.Errorf("ACCESS_KEY[%s] or ACCESS_SECRET[%s] is nil", key, secret)
		return nil, fmt.Errorf("accesskey or secret is nil")
	}

	/*
		deviceName := os.Getenv("HOSTNAME")
		if len(deviceName) == 0 {
			klog.Errorf("Need HOSTNAME ENV")
			return nil, fmt.Errorf("need HOSTNAME ENV")
		}
	*/

	h := &Mqtt3Handler{
		SubTopicHandlers: subTopicsHandlers,
	}

	clientID := fmt.Sprintf("%s@@@%s-sub", group, hostname)
	username := fmt.Sprintf("Signature|%s|%s", key, instance)
	passwd := util.GetSignature(clientID, secret)

	h.MqttSubClient = NewMqtt3Client(broker,
		port, clientID, username, passwd, true, true, h.Reconnect, h.LostConnectHandler)

	pubClientID := fmt.Sprintf("%s@@@%s-pub", group, hostname)
	pubUsername := fmt.Sprintf("Signature|%s|%s", key, instance)
	pubPasswd := util.GetSignature(pubClientID, secret)

	h.MqttDataClient = NewMqtt3Client(broker, port, pubClientID, pubUsername, pubPasswd, true, true, h.ReconnectNoSub, h.LostConnectHandler)

	ackClientID := fmt.Sprintf("%s@@@%s-ack", group, hostname)
	ackUsername := fmt.Sprintf("Signature|%s|%s", key, instance)
	ackPasswd := util.GetSignature(ackClientID, secret)

	h.MqttAckClient = NewMqtt3Client(broker, port, ackClientID, ackUsername, ackPasswd, true, true, h.ReconnectNoSub, h.LostConnectHandler)

	return h, nil
}

func (m *Mqtt3Handler) PublishData(ctx context.Context, topic string, qos byte, retained bool, object interface{}) error {
	return Mqtt3Send(m.MqttDataClient, topic, qos, retained, object)
}
func (m *Mqtt3Handler) PublishAck(ctx context.Context, topic string, qos byte, retained bool, object interface{}) error {
	return Mqtt3Send(m.MqttAckClient, topic, qos, retained, object)
}

var _ MessageHandler = &Mqtt3Handler{}

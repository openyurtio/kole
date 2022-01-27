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
	"encoding/json"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"k8s.io/klog/v2"
)

type SingleSubcribe struct {
	Topic   string
	Handler func(publish *paho.Publish)
	Option  paho.SubscribeOptions
}

func NewMqtt5Manager(
	ctx context.Context,
	keepAlive uint16,
	sessionExpiryInterval uint32,
	cleanStart bool,
	receiveMaximum uint16,
	connectRetryDelay time.Duration,
	packetTimeout time.Duration,
	server string,
	clientid string,
	subs []*SingleSubcribe,
) (*autopaho.ConnectionManager, error) {

	serverURL, err := url.Parse(server)
	if err != nil {
		klog.Errorf("Url parse %s error %v", server, err)
		return nil, err
	}

	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{serverURL},
		KeepAlive:         keepAlive,
		ConnectRetryDelay: connectRetryDelay,
		OnConnectError: func(err error) {
			klog.Errorf("Whilst attempting connection cliendid %s error: %s\n", clientid, err)
		},
		ClientConfig: paho.ClientConfig{
			ClientID: clientid,
			OnClientError: func(err error) {
				klog.Errorf("Server requested disconnect clientid %s: %s", clientid, err)
			},
			PacketTimeout: packetTimeout,
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					klog.Errorf("Server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					klog.Errorf("Server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}

	if len(subs) != 0 {
		router := paho.NewStandardRouter()
		subscribeOps := make(map[string]paho.SubscribeOptions)

		for _, single := range subs {
			router.RegisterHandler(single.Topic, single.Handler)
			subscribeOps[single.Topic] = single.Option
		}

		cliCfg.OnConnectionUp = func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			klog.V(5).Infof("Clientid %s mqtt connection up", clientid)
			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: subscribeOps,
			}); err != nil {
				klog.Errorf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
				return
			}
			klog.V(5).Infof("Clientid %s mqtt subscription made", clientid)
		}
		cliCfg.ClientConfig.Router = router
	} else {
		cliCfg.OnConnectionUp = func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			klog.V(5).Infof("Clientid %s mqtt connection up", clientid)
		}
	}

	cliCfg.SetConnectPacketConfigurator(func(connect *paho.Connect) *paho.Connect {
		connect.CleanStart = cleanStart
		connect.Properties = &paho.ConnectProperties{
			SessionExpiryInterval: &sessionExpiryInterval,
			ReceiveMaximum:        &receiveMaximum,
		}
		return connect
	})

	//
	// Connect to the broker
	//
	cm, err := autopaho.NewConnection(ctx, cliCfg)
	if err != nil {
		klog.Errorf("NewConnection clientid %s error %v", clientid, err)
		return nil, err
	}

	return cm, nil
}

func Mqtt5Send(ctx context.Context,
	cm *autopaho.ConnectionManager,
	topic string,
	qos byte,
	retained bool,
	object interface{}) error {

	data, err := json.Marshal(object)
	if err != nil {
		klog.Errorf("Mqtt5Send topic %s marshal error %v", topic, err)
		return err
	}

	// AwaitConnection will return immediately if connection is up; adding this call stops publication whilst
	// connection is unavailable.
	if err := cm.AwaitConnection(ctx); err != nil {
		// Should only happen when context is cancelled
		klog.Errorf("publisher done (AwaitConnection: %s)\n", err)
		return err
	}

	pr, err := cm.Publish(ctx, &paho.Publish{
		QoS:     qos,
		Topic:   topic,
		Payload: data,
		Retain:  retained,
	})
	if err != nil {
		klog.Errorf("error publishing: %s", err)
		return err
	} else if pr != nil && pr.ReasonCode != 0 && pr.ReasonCode != 16 { // 16 = Server received message but there are no subscribers
		klog.Errorf("reason code %d received\n", pr.ReasonCode)
		return err
	}
	klog.V(4).Infof("sent %s message: %s", topic, data)
	return nil
}

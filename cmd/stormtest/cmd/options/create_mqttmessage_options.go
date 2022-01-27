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
package options

import (
	"github.com/spf13/cobra"
)

type CreateMqttMessageFlags struct {
	*GlobalFlags
	MessageNum     int
	MqttBroker     string
	MqttBrokerPort int
	MqttGroup      string
	MqttInstance   string
	AccessKey      string
	AccessSecret   string
	Debug          bool
	Topic          string
	RandClientID   bool
}

func NewCreateMqttMessageFlags(g *GlobalFlags) *CreateMqttMessageFlags {
	return &CreateMqttMessageFlags{
		GlobalFlags:    g,
		MessageNum:     100000,
		MqttBroker:     "mqtt-cn-7mz2ietw201.mqtt.aliyuncs.com",
		MqttBrokerPort: 8883,
		MqttGroup:      "GID_TEST",
		MqttInstance:   "mqtt-cn-7mz2ietw201",
		Debug:          false,
		Topic:          "storm-rocket",
		RandClientID:   false,
	}
}

// AddFlags adds flags for a specific
func (f *CreateMqttMessageFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&f.MessageNum, "message-num", f.MessageNum, "the number of messages")
	cmd.Flags().StringVar(&f.MqttBroker, "mqtt-broker", f.MqttBroker, "the address of mqtt broker")
	cmd.Flags().StringVar(&f.Topic, "topic", f.Topic, "the topic of mqtt")
	cmd.Flags().IntVar(&f.MqttBrokerPort, "mqtt-broker-port", f.MqttBrokerPort, "the port of mqtt broker")
	cmd.Flags().StringVar(&f.MqttGroup, "mqtt-group", f.MqttGroup, "the mqtt group")
	cmd.Flags().StringVar(&f.MqttInstance, "mqtt-instance", f.MqttInstance, "mqtt instance name")

	cmd.Flags().StringVar(&f.AccessKey, "access-key", f.AccessKey, "access key")
	cmd.Flags().StringVar(&f.AccessSecret, "access-secret", f.AccessSecret, "access secret")
	cmd.Flags().BoolVar(&f.Debug, "debug", f.Debug, "is need debug")
	cmd.Flags().BoolVar(&f.RandClientID, "randclientid", f.RandClientID, "whether auto create id random")
}

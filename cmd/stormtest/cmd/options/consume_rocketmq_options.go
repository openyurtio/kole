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

type RocketMQFlags struct {
	*GlobalFlags
	SimulationsNums int
	MessageNum      int
	Endpoint        string
	Group           string
	Instance        string
	AccessKey       string
	AccessSecret    string
}

func NewConsumeRocketMQFlags(g *GlobalFlags) *RocketMQFlags {
	return &RocketMQFlags{
		GlobalFlags:     g,
		SimulationsNums: 1,
		MessageNum:      10,
		Endpoint:        "http://1255832284162844.mqrest.cn-beijing.aliyuncs.com",
		Group:           "GID_STROM",
		Instance:        "MQ_INST_1255832284162844_BX5q5eI3",
	}
}

// AddFlags adds flags for a specific
func (f *RocketMQFlags) AddFlags(cmd *cobra.Command) {
	// Here you will define your flags and configuration settings.
	cmd.Flags().IntVar(&f.SimulationsNums, "simulations-nums", f.SimulationsNums, "simulations-nums")
	cmd.Flags().IntVar(&f.MessageNum, "message-nums", f.MessageNum, "message-nums")
	cmd.Flags().StringVar(&f.Endpoint, "endpoint", f.Endpoint, "the endpoint of rocket broker")
	cmd.Flags().StringVar(&f.Group, "group", f.Group, "the rocket mq group")
	cmd.Flags().StringVar(&f.Instance, "instance", f.Instance, "rocketmq instance name")

	cmd.Flags().StringVar(&f.AccessKey, "access-key", f.AccessKey, "access key")
	cmd.Flags().StringVar(&f.AccessSecret, "access-secret", f.AccessSecret, "access secret")
}

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
package create

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/stormtest"
)

func NewCreateMessageCommand() *cobra.Command {

	var createMessageOptions *options.CreateMqttMessageFlags

	// createnodesCmd represents the createnodes command
	var createmessageCmd = &cobra.Command{
		Use:   "mqttmessage",
		Short: "Create a specific number of MQTT messages",
		Long:  "Create a specific number of MQTT messages",
		Run: func(cmd *cobra.Command, args []string) {

			klog.V(4).Infof("Stormtest create mqtt message config: %#v", *createMessageOptions)
			if createMessageOptions.Debug {
				mqtt.DEBUG = log.New(os.Stdout, "", 0)
				mqtt.ERROR = log.New(os.Stdout, "", 0)
			}
			if err := RunCreateMqttMessage(createMessageOptions); err != nil {
				klog.Fatal(err)
			}
		},
	}
	createMessageOptions = options.NewCreateMqttMessageFlags(&globalOptions)
	createMessageOptions.AddFlags(createmessageCmd)
	return createmessageCmd
}
func init() {
	subrootCmd.AddCommand(NewCreateMessageCommand())
}

func RunCreateMqttMessage(opts *options.CreateMqttMessageFlags) error {
	defer runtime.HandleCrash()

	lite, err := stormtest.NewCreateMqttMessage(opts)
	if err != nil {
		klog.Errorf("RunCreateMqttMessage error %v", err)
		return err
	}
	return lite.Run()
}

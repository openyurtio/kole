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
package consume

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/stormtest"
)

func NewConsumeMqttCommand() *cobra.Command {

	var consumeMqttOptions *options.ConsumeMqttFlags

	var consumeMqttCmd = &cobra.Command{
		Use:   "mqtt",
		Short: `mqtt message`,
		Long:  `mqtt message`,
		Run: func(cmd *cobra.Command, args []string) {

			klog.V(4).Infof("Stormtest consume mqtt message config: %#v", *consumeMqttOptions)
			if consumeMqttOptions.Debug {
				mqtt.DEBUG = log.New(os.Stdout, "", 0)
				mqtt.ERROR = log.New(os.Stdout, "", 0)
			}
			if err := RunMqtt(consumeMqttOptions); err != nil {
				klog.Fatal(err)
			}
		},
	}
	consumeMqttOptions = options.NewConsumeMqttFlags(&globalOptions)
	consumeMqttOptions.AddFlags(consumeMqttCmd)
	return consumeMqttCmd
}

func init() {
	subrootCmd.AddCommand(NewConsumeMqttCommand())
}

func RunMqtt(opts *options.ConsumeMqttFlags) error {
	defer runtime.HandleCrash()

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	group := sync.WaitGroup{}
	n := time.Now()
	for i := 0; i < opts.SimulationsNums; i++ {
		group.Add(1)
		go func(index int) {
			defer group.Done()
			lite, err := stormtest.NewConsumeMqtt(opts)
			if err != nil {
				klog.Errorf("NewConsumeMqtt error %v", err)
				return
			}
			lite.Run(ctx)
		}(i)
	}
	group.Wait()

	klog.Infof("SimulationsNums %d , Single Mqtt Message Num %d, use %d ms", opts.SimulationsNums,
		opts.MessageNum, time.Now().Sub(n).Milliseconds())

	return nil
}

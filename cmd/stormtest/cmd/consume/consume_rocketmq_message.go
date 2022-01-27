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
	"sync"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/stormtest"
)

func NewConsumemqCommand() *cobra.Command {

	var opts *options.RocketMQFlags

	var consumeMqCmd = &cobra.Command{
		Use:   "rocketmq",
		Short: (`rocketmq message`),
		Long:  (`rocketmq message`),
		Run: func(cmd *cobra.Command, args []string) {

			klog.V(4).Infof("Stormtest consume rocketmq message config: %#v", *opts)
			if err := RunMQ(opts); err != nil {
				klog.Fatal(err)
			}
		},
	}
	opts = options.NewConsumeRocketMQFlags(&globalOptions)
	opts.AddFlags(consumeMqCmd)
	return consumeMqCmd
}

func init() {
	subrootCmd.AddCommand(NewConsumemqCommand())
}

func RunMQ(opts *options.RocketMQFlags) error {
	defer runtime.HandleCrash()
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	group := sync.WaitGroup{}

	n := time.Now()
	for i := 0; i < opts.SimulationsNums; i++ {
		group.Add(1)
		go func(index int) {
			defer group.Done()
			lite, err := stormtest.NewConsumeRocketMQ(opts)
			if err != nil {
				klog.Errorf("NewConsumeRocketMQ error %v", err)
				return
			}
			lite.Run(ctx)
		}(i)
	}

	group.Wait()
	klog.Infof("SimulationsNums %d , Single Message Num %d, use %d ms", opts.SimulationsNums,
		opts.MessageNum, time.Now().Sub(n).Milliseconds())

	return nil
}

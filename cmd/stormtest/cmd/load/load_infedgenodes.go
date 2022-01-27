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
package load

import (
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/stormtest"
)

func NewLoadNodesCommand() *cobra.Command {

	var loadOptions *options.LoadInfEdgeNodeFlags

	// loadnodesCmd represents the loadnodes command
	var loadnodesCmd = &cobra.Command{
		Use:   "infedgenode",
		Short: "All computing load under a particular namespace infedgenodes instance need time (ms)",
		Long:  "All computing load under a particular namespace infedgenodes instance need time (ms)",
		Run: func(cmd *cobra.Command, args []string) {
			// run the kubelet
			klog.V(4).Infof("Stormtest load infedgenodes config: %#v", *loadOptions)
			if err := RunLoadInfEdgeNode(loadOptions); err != nil {
				klog.Fatal(err)
			}
		},
	}
	loadOptions = options.NewLoadInfEdgeNodeFlags(&globalOptions)
	loadOptions.AddFlags(loadnodesCmd)
	return loadnodesCmd
}
func init() {
	subrootCmd.AddCommand(NewLoadNodesCommand())
}

func RunLoadInfEdgeNode(config *options.LoadInfEdgeNodeFlags) error {
	defer runtime.HandleCrash()

	lite, err := stormtest.NewLoadInfEdgeNodes(config)
	if err != nil {
		klog.Errorf("NewLoadInfEdgeNodes error %v", err)
		return err
	}
	return lite.Run()
}

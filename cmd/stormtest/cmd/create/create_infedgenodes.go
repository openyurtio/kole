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
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/stormtest"
)

func NewCreateInfEdgeNodeCommand() *cobra.Command {
	var createOptions *options.CreateInfEdgeNodesFlags
	createOptions = options.NewCreateInfEdgeNodesFlags(&globalOptions)

	// createnodesCmd represents the createnodes command
	createnodesCmd := &cobra.Command{
		Use:   "infedgenode",
		Short: "In a particular namespace created under a certain number of infedgenodes. Lite. Openyurt. IO instance",
		Long:  "In a particular namespace created under a certain number of infedgenodes. Lite. Openyurt. IO instance",
		Run: func(cmd *cobra.Command, args []string) {

			// run the kubelet
			klog.V(4).Infof("Stormtest create infedgenodes config: %#v", *createOptions)
			if err := RunCreateInfEdgeNode(createOptions); err != nil {
				klog.Fatal(err)
			}
		},
	}

	createOptions.AddFlags(createnodesCmd)
	return createnodesCmd
}

func init() {
	subrootCmd.AddCommand(NewCreateInfEdgeNodeCommand())
}

func RunCreateInfEdgeNode(opts *options.CreateInfEdgeNodesFlags) error {
	defer runtime.HandleCrash()

	lite, err := stormtest.NewCreateInfEdgeNodes(opts)
	if err != nil {
		klog.Errorf("RunCreateInfEdgeNode error %v", err)
		return err
	}
	return lite.Run()
}

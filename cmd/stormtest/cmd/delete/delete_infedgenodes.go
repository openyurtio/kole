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
package delete

import (
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/stormtest"
)

func NewDeleteInfEdgeNodesCommand() *cobra.Command {

	var deleteOptions *options.DeleteInfEdgeNodesFlags

	// deletenodesCmd represents the deletenodes command
	var deletenodesCmd = &cobra.Command{
		Use:   "infedgenode",
		Long:  "Under the specific namespace to delete a certain number of infedgenodes instance",
		Short: "Under the specific namespace to delete a certain number of infedgenodes instance",
		Run: func(cmd *cobra.Command, args []string) {
			// run the kubelet
			klog.V(4).Infof("Stormtest delete infedgenodes config: %#v", *deleteOptions)
			if err := RunDeleteInfEdgeNode(deleteOptions); err != nil {
				klog.Fatal(err)
			}
		},
	}
	deleteOptions = options.NewDeleteInfEdgeNodesFlags(&globalOptions)
	deleteOptions.AddFlags(deletenodesCmd)
	return deletenodesCmd
}
func init() {
	subrootCmd.AddCommand(NewDeleteInfEdgeNodesCommand())
}

func RunDeleteInfEdgeNode(config *options.DeleteInfEdgeNodesFlags) error {
	defer runtime.HandleCrash()

	lite, err := stormtest.NewDeleteInfEdgeNodes(config)
	if err != nil {
		klog.Errorf("NewDeleteInfEdgeNodes error %v", err)
		return err
	}
	return lite.Run()

}

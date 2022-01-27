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

func NewCreateSummariesCommand() *cobra.Command {

	var createSumOptions *options.CreateSummariesFlags

	// createnodesCmd represents the createnodes command
	var createsumarisCmd = &cobra.Command{
		Use:   "summaries",
		Short: "Create a number of Summaris instances in a specific namespace",
		Long:  "Create a number of Summaris instances in a specific namespace",
		Run: func(cmd *cobra.Command, args []string) {

			klog.V(4).Infof("Stormtest create summaries config: %#v", *createSumOptions)
			if err := RunCreateSummaris(createSumOptions); err != nil {
				klog.Fatal(err)
			}
		},
	}
	createSumOptions = options.NewCreateSummariesFlags(&globalOptions)
	createSumOptions.AddFlags(createsumarisCmd)
	return createsumarisCmd
}

func init() {
	subrootCmd.AddCommand(NewCreateSummariesCommand())
}

func RunCreateSummaris(opts *options.CreateSummariesFlags) error {
	defer runtime.HandleCrash()

	lite, err := stormtest.NewCreateSummaries(opts)
	if err != nil {
		klog.Errorf("RunCreateSummaris error %v", err)
		return err
	}
	return lite.Run()
}

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
package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/stormtest"
)

func NewCompressCommand() *cobra.Command {

	var compressOptions *options.CompressionFlags

	var compresscrsCmd = &cobra.Command{
		Use:   "compression",
		Short: "Test with different compression algorithms",
		Long:  "Test with different compression algorithms",
		Run: func(cmd *cobra.Command, args []string) {
			klog.V(4).Infof("Stormtest compresscrs config: %#v", *compressOptions)
			if err := RunCompress(compressOptions); err != nil {
				klog.Fatal(err)
			}
		},
	}
	compressOptions = options.NewCompressionFlags(&globalOptions)
	compressOptions.AddFlags(compresscrsCmd)
	return compresscrsCmd
}
func init() {
	rootCmd.AddCommand(NewCompressCommand())
}

func RunCompress(opt *options.CompressionFlags) error {
	defer runtime.HandleCrash()

	test, err := stormtest.NewCompression(opt)
	if err != nil {
		klog.Errorf("NewCompression error %v", err)
		return err
	}
	return test.Run()
}

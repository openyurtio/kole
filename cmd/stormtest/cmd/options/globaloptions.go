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
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type GlobalFlags struct {
	KubeConfig string
}

// ValidateGlobalFlags
func ValidateGlobalFlags(f *GlobalFlags) error {
	// ensure that nobody sets DynamicConfigDir if the dynamic config feature gate is turned off

	return nil
}

// ValidateCreateNodesFlags
func ValidateCreateNodesFlags(f *CreateInfEdgeNodesFlags) error {
	// ensure that nobody sets DynamicConfigDir if the dynamic config feature gate is turned off

	return nil
}

func NewGlobalFlags() *GlobalFlags {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &GlobalFlags{
		KubeConfig: filepath.Join(home, ".kube/config"),
	}
}

// AddFlags adds flags for a specific
func (f *GlobalFlags) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&f.KubeConfig, "kube-config", f.KubeConfig, "config file (default is $HOME/.kube/config)")
	//cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

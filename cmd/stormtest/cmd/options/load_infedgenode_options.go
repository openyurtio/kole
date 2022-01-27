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

type LoadInfEdgeNodeFlags struct {
	*GlobalFlags
	NS       string
	PatchNum int
}

func NewLoadInfEdgeNodeFlags(g *GlobalFlags) *LoadInfEdgeNodeFlags {
	return &LoadInfEdgeNodeFlags{
		GlobalFlags: g,
		NS:          "summarystorm",
		PatchNum:    500,
	}
}

// AddFlags adds flags for a specific
func (f *LoadInfEdgeNodeFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&f.PatchNum, "patch-num", f.PatchNum, "the list chunk max num")
}

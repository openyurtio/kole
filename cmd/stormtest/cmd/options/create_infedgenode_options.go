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

type CreateInfEdgeNodesFlags struct {
	*GlobalFlags
	NodeNum     int
	BatchNum    int
	NS          string
	IsSmallSize bool
}

func NewCreateInfEdgeNodesFlags(g *GlobalFlags) *CreateInfEdgeNodesFlags {
	return &CreateInfEdgeNodesFlags{
		GlobalFlags: g,
		NodeNum:     100,
		BatchNum:    10,
		NS:          "summarystorm",
		IsSmallSize: true,
	}
}

// AddFlags adds flags for a specific
func (f *CreateInfEdgeNodesFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&f.NodeNum, "node-num", f.NodeNum, "the number of nodes")
	cmd.Flags().IntVar(&f.BatchNum, "batch-num", f.BatchNum, "the number of batch nums")
	cmd.Flags().BoolVar(&f.IsSmallSize, "is-small-size", f.IsSmallSize,
		"Whether to use small size to create,  when false it use 8k size to create, when true, it use 1 k size to create")
}

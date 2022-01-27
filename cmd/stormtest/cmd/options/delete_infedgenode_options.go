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

type DeleteInfEdgeNodesFlags struct {
	*GlobalFlags
	NS             string
	DeletePatchNum int
}

func NewDeleteInfEdgeNodesFlags(g *GlobalFlags) *DeleteInfEdgeNodesFlags {
	return &DeleteInfEdgeNodesFlags{
		GlobalFlags:    g,
		NS:             "summarystorm",
		DeletePatchNum: 100,
	}
}

// AddFlags adds flags for a specific
func (f *DeleteInfEdgeNodesFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&f.DeletePatchNum, "delete-patch-num", f.DeletePatchNum, "delete the patch nums of nodes")
}

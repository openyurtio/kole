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

	"github.com/spf13/cobra"

	"github.com/openyurtio/kole/pkg/controller"
)

type CompressionFlags struct {
	*GlobalFlags
	SrcNs       string
	DstNs       string
	Algthm      string
	OnlySaveCrs bool
}

var supportAlgthms map[string]controller.DataProcesser

func init() {
	supportAlgthms = make(map[string]controller.DataProcesser)
	supportAlgthms["gzip"] = &controller.Gzip{}
	supportAlgthms["lzw"] = &controller.Lzw{}
	supportAlgthms["flate"] = &controller.Flate{}
	supportAlgthms["lz4"] = &controller.Lz4{}
	supportAlgthms["snappy"] = &controller.Snappy{}
}

func SupprtAlgthms() []string {
	alg := make([]string, 0, 10)
	for k, _ := range supportAlgthms {
		alg = append(alg, k)
	}
	return alg
}

func AlgthmFactory(algName string) (controller.DataProcesser, error) {

	p, ok := supportAlgthms[algName]
	if ok {
		return p, nil
	}
	return nil, fmt.Errorf("Do not supprt %s algthm", algName)
}

func NewCompressionFlags(g *GlobalFlags) *CompressionFlags {
	return &CompressionFlags{
		GlobalFlags: g,
		SrcNs:       "jinchen",
		DstNs:       "compresstest",
		Algthm:      "gzip",
		OnlySaveCrs: true,
	}
}

func (f *CompressionFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.Algthm, "alg-name", f.Algthm, fmt.Sprintf("The name of compress algorithm, such as %++v", SupprtAlgthms()))
	cmd.Flags().StringVar(&f.DstNs, "dst-ns", f.DstNs, "The destination namespace")
	cmd.Flags().StringVar(&f.SrcNs, "src-ns", f.SrcNs, "The source namespace")
	cmd.Flags().BoolVar(&f.OnlySaveCrs, "only-save", f.OnlySaveCrs, "Whether to save cr from src ns to dst ns,true means you only save cr from src ns to dst ns")
}

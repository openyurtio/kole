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
package stormtest

import (
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
	"github.com/openyurtio/kole/pkg/controller"
)

func TestSummris(t *testing.T) {
	tt := "0000000000"
	var s string
	var ss string
	for i := 0; i < 1000; i++ {
		s = fmt.Sprintf("%s%s", s, tt)
	}
	for i := 0; i < 1000; i++ {
		ss = fmt.Sprintf("%s%s", ss, s)
	}
	t.Logf("The length of data is %d", len(ss))

	sum := &v1alpha1.Summary{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "sdf",
			Name:      "sdf",
		},
		Data:  []byte(ss),
		Index: 1,
	}

	result, err := yaml.Marshal(sum)
	if err != nil {
		t.Errorf("Marshal error %v", err)
		return
	}

	t.Logf("The length of data after marshal is %d", len(result))

	g := &controller.Gzip{}
	compressData, err := g.Compress(result)
	if err != nil {
		t.Errorf("Compress error %v", err)
		return
	}
	t.Logf("The length of compressed data is %d", len(compressData))

	sum = &v1alpha1.Summary{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "sdf",
			Name:      "sdf",
		},
		Data:  compressData,
		Index: 1,
	}

	result, err = yaml.Marshal(sum)
	if err != nil {
		t.Errorf("Marshal error %v", err)
		return
	}

	t.Logf("The length of data after compress marshal is %d", len(result))
}

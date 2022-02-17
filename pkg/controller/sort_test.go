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
package controller

import (
	"sort"
	"testing"

	"github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
)

func TestBySummary(t *testing.T) {
	cases := []struct {
		Name string
		List []v1alpha1.Summary
	}{
		{
			"first",
			[]v1alpha1.Summary{
				v1alpha1.Summary{
					Data:  []byte("second"),
					Index: 2,
				},
				v1alpha1.Summary{
					Data:  []byte("first"),
					Index: 1,
				},
				v1alpha1.Summary{
					Data:  []byte("-first"),
					Index: -1,
				},
				v1alpha1.Summary{
					Data:  []byte("hahh 500"),
					Index: 50,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			sort.Stable(BySummary(c.List))
			for _, l := range c.List {
				t.Logf(" data %s index %d", l.Data, l.Index)
			}
		})
	}
}

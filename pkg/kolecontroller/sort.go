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
package kolecontroller

import "github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"

// BySummary implements sort.Interface to allow ControllerRevisions to be sorted by Revision.
type BySummary []v1alpha1.Summary

func (br BySummary) Len() int {
	return len(br)
}

// Less breaks ties first by creation timestamp, then by name
func (br BySummary) Less(i, j int) bool {
	return br[i].Index < br[j].Index
}

func (br BySummary) Swap(i, j int) {
	br[i], br[j] = br[j], br[i]
}

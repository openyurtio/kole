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

import (
	"sync"

	"github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
)

type QueryNodeStatusCache struct {
	*sync.RWMutex
	NameToStatus map[string]*v1alpha1.QueryNodeStatus
}

func (c *QueryNodeStatusCache) Reset(nameToStatus map[string]*v1alpha1.QueryNodeStatus) {
	c.Lock()
	c.NameToStatus = nameToStatus
	c.Unlock()
}

func (c *QueryNodeStatusCache) GetNodeStatus(nodeName string) *v1alpha1.QueryNodeStatus {
	var s *v1alpha1.QueryNodeStatus
	c.RLock()
	s = c.NameToStatus[nodeName]
	c.RUnlock()
	return s
}

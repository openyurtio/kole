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
	"sync"

	"github.com/openyurtio/kole/pkg/data"
)

type ObserverdPodsCache struct {
	*sync.RWMutex
	// nodeName / pod key / Pod
	Cache map[string]map[string]*data.HeartBeatPod
}

func (c *ObserverdPodsCache) SafeSetHeartBeat(hb *data.HeartBeat) {
	c.Lock()
	_, ok := c.Cache[hb.Name]
	if !ok {
		c.Cache[hb.Name] = make(map[string]*data.HeartBeatPod)
	}
	for _, hbp := range hb.Pods {
		c.Cache[hb.Name][hbp.Key()] = &data.HeartBeatPod{
			Hash:      hbp.Hash,
			Name:      hbp.Name,
			NameSpace: hbp.NameSpace,
			Status:    hbp.Status,
		}
	}
	c.Unlock()
}
func (c *ObserverdPodsCache) ReadRange(f func(nodeName string, hbPodList map[string]*data.HeartBeatPod)) {
	c.RLock()
	for nodeName, _ := range c.Cache {
		f(nodeName, c.Cache[nodeName])
	}
	c.RUnlock()
}

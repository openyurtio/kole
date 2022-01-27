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

	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/pkg/data"
)

type FilterInfo struct {
	SeqNum    uint64
	TimeStamp int64
}

type HeatBeatFilter struct {
	*sync.Mutex
	Filter map[string]*FilterInfo
}

func (c *HeatBeatFilter) SetHeatBeat(hb *data.HeatBeat) bool {
	c.Lock()
	setSuccess := true

	old, ok := c.Filter[hb.Name]
	if ok {
		if hb.SeqNum < old.SeqNum {
			klog.Warningf("Receive HeatBeat Node %s Seq %v is less then cache seq %v, do nothing", hb.Name, hb.SeqNum, old.SeqNum)
			setSuccess = false
		} else if hb.SeqNum == old.SeqNum && hb.TimeStamp <= old.TimeStamp {
			klog.V(4).Infof("Receive HeatBeat Node %s Seq %v equal cache seq, but timestamp[%v] is less or equal than cache[%v], do nothing", hb.Name, hb.SeqNum,
				hb.TimeStamp, old.TimeStamp)
			setSuccess = false
		}
	}

	if setSuccess {
		c.Filter[hb.Name] = &FilterInfo{
			SeqNum:    hb.SeqNum,
			TimeStamp: hb.TimeStamp,
		}
	}
	c.Unlock()
	return setSuccess
}

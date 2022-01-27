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
package litekubelet

import (
	"sync"

	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/pkg/data"
)

var localPods map[string]*data.Pod
var localPodsLock sync.Mutex

func init() {
	localPods = make(map[string]*data.Pod)
}

func SyncLocalPod(p *data.Pod) {
	localPodsLock.Lock()
	defer localPodsLock.Unlock()
	if p.DeleteTimeStamp != nil {
		klog.Infof("Delete pod %s", p.Key())
		delete(localPods, p.Key())
	} else {
		localPods[p.Key()] = p
	}
}

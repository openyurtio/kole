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
package cache

import (
	"sync"
	"time"

	"k8s.io/klog/v2"
)

type TimeoutCache struct {
	cache sync.Map
	cond  *sync.Cond
}

type cacheData struct {
	value interface{}
	t     time.Time
}

var timecache *TimeoutCache

func init() {
	timecache = NewTimeoutCache(time.Second*5, time.Minute*120)
}

func GetDefaultTimeoutCache() *TimeoutCache {
	return timecache
}

func NewTimeoutCache(period, timeout time.Duration) *TimeoutCache {
	t := &TimeoutCache{
		cache: sync.Map{},
		cond:  sync.NewCond(&sync.Mutex{}),
	}
	go t.run(period, timeout)
	return t
}

// 每个周期清理缓存中过期的key, 每个周期广播一次，防止缓存里的数据过多，减少内存占用
func (t *TimeoutCache) run(period, timeout time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n := time.Now()
			t.cache.Range(func(key, value interface{}) bool {
				if n.Sub(value.(*cacheData).t) > timeout {
					klog.Infof("key %s data %d now %d exceed %d s", key, value.(*cacheData).t.Unix(), n.Unix(), timeout.Milliseconds()/1000)
					t.cache.Delete(key)
				}
				return true
			})
			t.cond.Broadcast()
		}
	}
}

// Pop retrun false , when timeout
// 获取key 的值， 如果有返回对应的value ,并返回true,同时从缓存里delete 掉
// 如果没有，会一直wait, 直到获得对应的key 和value
// 若一直到timeout 的时间，没有获得数据，返回false
func (t *TimeoutCache) PopWait(key string, timeout time.Duration) (interface{}, bool) {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			return nil, false
		default:
		}
		v, ok := t.cache.Load(key)
		if ok {
			t.cache.Delete(key)
			return v.(*cacheData).value, true
		}
		t.cond.L.Lock()
		t.cond.Wait()
		t.cond.L.Unlock()
	}
}

// 向超时缓存里存放key 和value, 一旦存放，则广播给哪些PopWait 中的查询，通知他们重新查询（sync.cond.Broadcast()）
func (t *TimeoutCache) Set(key string, value interface{}) {
	d := &cacheData{
		value: value,
		t:     time.Now(),
	}
	t.cache.Store(key, d)
	t.cond.Broadcast()
}

// Get retrun false , if not exist
// 获取超时缓存里的数据， 有则返回对应的值，和true, 没有则返回false
func (t *TimeoutCache) Get(key string) (interface{}, bool) {
	v, ok := t.cache.Load(key)
	if ok {
		return v.(*cacheData).value, true
	}
	return nil, false
}

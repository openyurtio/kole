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
	"testing"
	"time"
)

func TestSub(t *testing.T) {
	n := time.Now()
	timeout := time.Second * 2
	time.Sleep(time.Second * 2)
	nn := time.Now()
	t.Logf("jiange %d", nn.Sub(n))
	t.Logf("jiange %d", nn.Sub(n).Milliseconds())
	t.Logf("timeout %d ms", timeout.Milliseconds())
	t.Logf("timeout %d", timeout)

	tn := time.Now().Unix()
	time.Sleep(time.Second * 2)
	tnn := time.Now().Unix()
	t.Logf("sub %d", tnn-tn)
}

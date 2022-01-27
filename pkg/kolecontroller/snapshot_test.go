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
	"testing"
)

func TestBytesBuffer(t *testing.T) {
	/*
		buffer := bytes.NewBuffer([]byte("zhangjie hello world!!!!a"))
			length := 4
			for {
				cache := make([]byte, length)
				n, err := buffer.Read(cache)
				if err != nil {
					t.Errorf("buffer read error %v ", err)
					break
				}
				if n < length {
					t.Logf("laster read , len %v, cache:[%s]", n, string(cache))
					break
				} else {
					t.Logf("read(%d) to cache:%s", n, string(cache))
				}
			}

		for {
			data := buffer.Next(4)
			if len(data) == 0 {
				break
			}
			t.Logf("read cache:%s", string(data))
		}
	*/
}

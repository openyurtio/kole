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
	"bytes"
	"fmt"
	"testing"
)

func TestGzip_Compress(t *testing.T) {
	context := []byte("List all InfEdgeNode in ns[summarystorm] error The provided continue parameter is too old to display a consistent list result. You can start a new list without the continue parameter, or use the continue token in this response to retrieve the remainder of the results. Continuing with the provided token results in an inconsistent list - objects that were created, modified, or deleted between the time the first chunk was returned and now may show up in the list.")
	var (
		gcmp       Gzip
		cmp, uncmp []byte
		err        error
	)
	if cmp, err = gcmp.Compress(context); err != nil {
		t.Errorf("compress err: %v", err)
	}
	fmt.Printf("the length of the origin date is :%d, the length after compress is :%d\n", len(context), len(cmp))
	//fmt.Printf("the context of the origin date is :%s\n, the length after compress is :%s\n", data, cmp)
	if uncmp, err = gcmp.UnCompress(cmp); err != nil {
		t.Errorf("uncompress err: %v", err)
	} else if bytes.Compare(context, uncmp) != 0 {
		t.Errorf("uncompressed data are not the same as original")
	}
}

func TestLzw_Compress(t *testing.T) {
	context := []byte("List all InfEdgeNode in ns[summarystorm] error The provided continue parameter is too old to display a consistent list result. You can start a new list without the continue parameter, or use the continue token in this response to retrieve the remainder of the results. Continuing with the provided token results in an inconsistent list - objects that were created, modified, or deleted between the time the first chunk was returned and now may show up in the list.")
	var (
		Lcmp       Lzw
		cmp, uncmp []byte
		err        error
	)
	if cmp, err = Lcmp.Compress(context); err != nil {
		t.Errorf("compress err: %v", err)
	}
	fmt.Printf("the length of the origin date is :%d, the length after compress is :%d\n", len(context), len(cmp))
	if uncmp, err = Lcmp.UnCompress(cmp); err != nil {
		t.Errorf("uncompress err: %v", err)
	} else if bytes.Compare(context, uncmp) != 0 {
		t.Errorf("uncompressed data are not the same as original")
	}
}

func TestFlate_Compress(t *testing.T) {
	context := []byte("List all InfEdgeNode in ns[summarystorm] error The provided continue parameter is too old to display a consistent list result. You can start a new list without the continue parameter, or use the continue token in this response to retrieve the remainder of the results. Continuing with the provided token results in an inconsistent list - objects that were created, modified, or deleted between the time the first chunk was returned and now may show up in the list.")
	var (
		Fcmp       Flate
		cmp, uncmp []byte
		err        error
	)
	if cmp, err = Fcmp.Compress(context); err != nil {
		t.Errorf("compress err: %v", err)
	}
	fmt.Printf("the length of the origin date is :%d, the length after compress is :%d\n", len(context), len(cmp))
	if uncmp, err = Fcmp.UnCompress(cmp); err != nil {
		t.Errorf("uncompress err: %v", err)
	} else if bytes.Compare(context, uncmp) != 0 {
		t.Errorf("uncompressed data are not the same as original")
	}
}

func TestLz4Compress_Compress(t *testing.T) {
	context := []byte("List all InfEdgeNode in ns[summarystorm] error The provided continue parameter is too old to display a consistent list result. You can start a new list without the continue parameter, or use the continue token in this response to retrieve the remainder of the results. Continuing with the provided token results in an inconsistent list - objects that were created, modified, or deleted between the time the first chunk was returned and now may show up in the list.")
	var (
		Lcmp       Lz4
		cmp, uncmp []byte
		err        error
	)
	if cmp, err = Lcmp.Compress(context); err != nil {
		t.Errorf("compress err: %v", err)
	}
	fmt.Printf("the length of the origin date is :%d, the length after compress is :%d\n", len(context), len(cmp))
	if uncmp, err = Lcmp.UnCompress(cmp); err != nil {
		t.Errorf("uncompress err: %v", err)
	} else if bytes.Compare(context, uncmp) != 0 {
		t.Errorf("uncompressed data are not the same as original")
	}
}

func TestSnappyCompress_Compress(t *testing.T) {
	context := []byte("List all InfEdgeNode in ns[summarystorm] error The provided continue parameter is too old to display a consistent list result. You can start a new list without the continue parameter, or use the continue token in this response to retrieve the remainder of the results. Continuing with the provided token results in an inconsistent list - objects that were created, modified, or deleted between the time the first chunk was returned and now may show up in the list.")
	var (
		Scmp       Snappy
		cmp, uncmp []byte
		err        error
	)
	if cmp, err = Scmp.Compress(context); err != nil {
		t.Errorf("compress err: %v", err)
	}
	fmt.Printf("the length of the origin date is :%d, the length after compress is :%d\n", len(context), len(cmp))
	if uncmp, err = Scmp.Compress(cmp); err != nil {
		t.Errorf("uncompress err: %v", err)
	}
	if uncmp, err = Scmp.UnCompress(cmp); err != nil {
		t.Errorf("uncompress err: %v", err)
	} else if bytes.Compare(context, uncmp) != 0 {
		t.Errorf("uncompressed data are not the same as original")
	}
}

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
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"io"
	"io/ioutil"
	"time"

	"github.com/golang/snappy"
	"github.com/pierrec/lz4"
	"k8s.io/klog/v2"
)

type DataProcesser interface {
	Compress(data []byte) ([]byte, error)
	UnCompress(data []byte) ([]byte, error)
}

type DataProcessNothing struct {
}

func (d *DataProcessNothing) Compress(data []byte) ([]byte, error) {
	return data, nil
}

func (d *DataProcessNothing) UnCompress(data []byte) ([]byte, error) {
	return data, nil
}

var _ DataProcesser = &DataProcessNothing{}

//Gzip Compress/EnCompress
type Gzip struct {
}

func (g *Gzip) Compress(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	//NewWriter方法初始化
	//now := time.Now()
	writer := gzip.NewWriter(&buffer)
	//使用压缩等级
	//writer, _ := gzip.NewWriterLevel(&buffer, 5)
	//最开始一直把关闭放在了defer里面，但是这样压缩后数据为空
	//查阅相关资料才知道：gzip压缩的过程中，Write之后一定要及时Close，不能defer，这样才能flush，否则得不到任何数据
	/*defer func() {
		err := writer.Close()
		if err != nil {
			logrus.Info("Gzip压缩失败", err.Error())
		}
	}()*/
	_, err := writer.Write(data)
	if err != nil {
		//write 出现err后要不要Close()
		//err = writer.Close()
		klog.Warning("Gzip Compress fail:", err)
		return nil, err
	}
	if err = writer.Close(); err != nil {
		klog.Warning("Close the Gzip Object fail:", err)
		return nil, err
	}
	//needTime := time.Now().Sub(now).Milliseconds()
	//klog.Info("implement Gzip Compress, need time:", needTime)
	return buffer.Bytes(), nil
}

func (g *Gzip) UnCompress(data []byte) ([]byte, error) {
	//now := time.Now()
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = reader.Close(); err != nil {
			println("Gzip EnCompress fail:", err.Error())
		}
	}()
	//needTime := time.Now().Sub(now).Milliseconds()
	//klog.Info("implement Gzip Encompress, need time:", needTime)
	return ioutil.ReadAll(reader)
}

//Lzw Compress/UnCompress
type Lzw struct {
}

func (L *Lzw) Compress(data []byte) ([]byte, error) {
	//now := time.Now()
	buf := bytes.NewBuffer(nil)
	w := lzw.NewWriter(buf, lzw.LSB, 8)
	_, err := w.Write(data)
	if err != nil {
		klog.Warning("Lzw Compress fail:", err)
		return nil, err
	}
	if err := w.Close(); err != nil {
		klog.Warning("Close the Lzw Compress Object fail :", err)
		return nil, err
	}
	//needTime := time.Now().Sub(now).Milliseconds()
	//klog.Info("implement Lzw Compress, need time:", needTime)
	return buf.Bytes(), nil
}

func (L *Lzw) UnCompress(data []byte) ([]byte, error) {
	now := time.Now()
	buf := bytes.NewBuffer(data)
	r := lzw.NewReader(buf, lzw.LSB, 8)
	//reader, err := lzw.NewReader(bytes.NewReader(data))

	defer func() {
		if err := r.Close(); err != nil {
			println("Close the Lzw EnCompress Object fail:", err.Error())
		}
	}()
	needTime := time.Now().Sub(now).Milliseconds()
	klog.Info("implement Lzw Encompress, need time:", needTime)
	return ioutil.ReadAll(r)
}

// go标准库中提供了一个对Bzip2压缩包进行读取的操作/解压操作
// 但是并没有提供进行bzip2压缩操作

//Flate Compress/EnCompress
type Flate struct {
}

func (f *Flate) Compress(data []byte) ([]byte, error) {
	now := time.Now()
	// 一个缓存区压缩的内容
	buf := bytes.NewBuffer(nil)

	// 创建一个flate.Writer
	flateWrite, err := flate.NewWriter(buf, flate.BestCompression)
	if err != nil {
		klog.Fatalln(err)
		return nil, err
	}
	defer flateWrite.Close()
	if _, err := flateWrite.Write(data); err != nil {
		klog.Info("Write data that will be compressed fail ", err)
		return nil, err
	}
	flateWrite.Flush()
	needTime := time.Now().Sub(now).Milliseconds()
	klog.Info("implement Flate Encompress, need time:", needTime)
	return buf.Bytes(), nil
}

func (f *Flate) UnCompress(data []byte) ([]byte, error) {
	now := time.Now()

	buf := bytes.NewBuffer(data)
	flateReader := flate.NewReader(buf)

	defer flateReader.Close()
	var rb, err = ioutil.ReadAll(flateReader)
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			klog.Fatalf(
				"Err %v, read %v", err, rb)
		}
	}
	needTime := time.Now().Sub(now).Milliseconds()
	klog.Info("implement Flate Encompress, need time: ", needTime)
	return rb, nil
}

//SnappyCompress realize the compress/uncompresss algorithm
type Snappy struct {
}

func (Snp *Snappy) Compress(data []byte) ([]byte, error) {
	cmp := snappy.Encode(nil, data)
	return cmp, nil
}

func (Snp *Snappy) UnCompress(data []byte) ([]byte, error) {
	uncmp, err := snappy.Decode(nil, data)
	if err != nil {
		klog.Errorf("snappy uncompress fail :%v", err)
		return nil, err
	}
	return uncmp, nil
}

//lz4 compress algorithm
type Lz4 struct {
}

func (lz *Lz4) Compress(data []byte) ([]byte, error) {
	cmp := make([]byte, len(data))
	//ht := make([]int, 64<<10)
	l, err := lz4.CompressBlock(data, cmp, nil)
	if err != nil {
		klog.Errorf("Lz4 algorithm compress fail %v", err)
		return nil, err
	}
	return cmp[:l], nil
}

func (lz *Lz4) UnCompress(data []byte) ([]byte, error) {
	uncmp := make([]byte, 10*len(data))
	l, err := lz4.UncompressBlock(data, uncmp)
	if err != nil {
		klog.Errorf("Lz4 algorithm uncompress fail %v", err)
		return nil, err
	}
	return uncmp[:l], nil
}

//nocompress
//type Nonecompress struct {
//}
//
//func (N *Nonecompress) Compress(data []byte) ([]byte, error){
//	return nil, nil
//}
//func (N *Nonecompress) UnCompress(data []byte) ([]byte, error){
//	return nil, nil
//}

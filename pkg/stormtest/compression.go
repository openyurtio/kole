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
package stormtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
	"github.com/openyurtio/kole/pkg/client/clientset/versioned"
	"github.com/openyurtio/kole/pkg/controller"
	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/util"
)

type Compression struct {
	KubeConfig string
	DstNs      string
	SrcNs      string
	Algthm     string
	OlSV       bool
	liteClient versioned.Interface
}

func NewCompression(config *options.CompressionFlags) (*Compression, error) {

	c, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		return nil, err
	}
	c.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(10000, 10000)
	// set rate limit

	crdclient, err := versioned.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	t := &Compression{
		KubeConfig: config.KubeConfig,
		SrcNs:      config.SrcNs,
		DstNs:      config.DstNs,
		Algthm:     config.Algthm,
		OlSV:       config.OnlySaveCrs,
		liteClient: crdclient,
	}
	return t, nil
}

//CreateSummariesByCache use cache context to create Cr in DstNs namespace
func (n *Compression) CreateSummariesByCache(dstNs, flag string, summariesCache []byte) (int, error) {
	klog.V(4).Infof("Start to create all summary cr...")
	bf := bytes.NewBuffer(summariesCache)
	totalNum := 0
	bufferLen := util.SNAPSHOT_MAX_BUFFER_LEN
	group := sync.WaitGroup{}
	for i := 0; ; i++ {
		data := bf.Next(bufferLen)
		if len(data) == 0 {
			break
		}

		//klog.Infof("Length of data is %d", len(data))
		sum := &v1alpha1.Summary{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: dstNs,
				Name:      fmt.Sprintf("%s-%d", flag, i),
			},
			Data:  data,
			Index: i,
		}

		totalNum++
		group.Add(1)
		go func(s *v1alpha1.Summary) {
			defer group.Done()
			for j := 0; j < 3; j++ {
				if _, err := n.liteClient.LiteV1alpha1().Summaries(s.Namespace).Create(context.Background(), s, metav1.CreateOptions{}); err != nil {
					klog.Errorf("create summary [%s][%s] error %v ", s.GetNamespace(), s.GetName(), err)
					time.Sleep(time.Second)
				} else {
					klog.V(4).Infof("create summary [%s][%s] successful", s.GetNamespace(), s.GetName())
					break
				}
			}
		}(sum)

	}
	group.Wait()
	klog.V(4).Infof("Save %d summares to ns %s successful", totalNum, n.DstNs)
	return totalNum, nil
}

//LoadSummariesByNs load the cr from DstNs namespace to cache
func (n *Compression) LoadSummariesByNs(ns string) ([]byte, error) {
	var timeoutS int64 = 60
	var continueStr string
	var max int64 = 500
	var total int

	cacheData := make([]byte, 0, 1024*1000*1000)
	allSummaris := make([]v1alpha1.Summary, 0, 1024)
	for {
		localSummaries, err := n.liteClient.LiteV1alpha1().Summaries(ns).List(context.Background(), metav1.ListOptions{
			TimeoutSeconds: &timeoutS,
			Limit:          max,
			Continue:       continueStr,
		})

		if err != nil {
			klog.Errorf("List all summarys in ns[%s] error %v", ns, err)
			return nil, err
		}
		getLen := len(localSummaries.Items)
		total = total + getLen
		if getLen == 0 {
			break
		}

		// TODO 可以使用bytes.buffer
		for i, _ := range localSummaries.Items {
			allSummaris = append(allSummaris, localSummaries.Items[i])
		}

		//klog.Infof("Current get list node chunk nums %d ", len(listedNodes.Items))
		continueStr = localSummaries.GetContinue()

		if len(continueStr) == 0 {
			klog.V(4).Infof("continueStr is null , break ")
			break
		}

		//klog.Infof("Current get all node chunk nums %d,  continueStr %s ", allNum, continueStr)
		if getLen < int(max) {
			break
		}

	}

	//sort by index
	sort.Stable(controller.BySummary(allSummaris))
	for i, _ := range allSummaris {
		cacheData = append(cacheData, allSummaris[i].Data...)
	}

	if total == 0 {
		klog.V(4).Infof("Can not get any summary cr")
		return nil, nil
	}
	klog.V(4).Infof("Get %d crs from ns %s", total, ns)
	return cacheData, nil
}

//delete cr in DstNs namespace
func (n *Compression) DeleteSummariesByNS(ns string) error {
	var timeoutS int64 = 60
	var continueStr string
	var max int64 = 500
	needDeleteList := make([]string, 0, 500)

	for {
		crs, err := n.liteClient.LiteV1alpha1().Summaries(ns).List(context.Background(), metav1.ListOptions{
			TimeoutSeconds: &timeoutS,
			Limit:          max,
			Continue:       continueStr,
		})
		if err != nil {
			klog.Errorf("List summaries cr from ns %s error %v", ns, err)
			return err
		}
		continueStr = crs.GetContinue()
		for _, single := range crs.Items {

			needDeleteList = append(needDeleteList, single.Name)
			//if n.DeleteNum > 0 && totalNum >= n.DeleteNum {
			//	deleteFinished = true
			//	break
			//}
		}
		//if len(crs.Items) < int(max) || deleteFinished {
		if len(crs.Items) < int(max) {
			break
		}
	}
	for _, sn := range needDeleteList {
		err := n.liteClient.LiteV1alpha1().Summaries(ns).Delete(context.Background(), sn, metav1.DeleteOptions{})
		if err != nil {
			klog.Errorf("Delete InfEdgeNode[%s][%s] error %v", ns, sn, err)
			return err
		}
	}
	klog.V(4).Infof("Delete all summaries cr in %s namespace success", ns)
	return nil
}

// ProcessHeartBeatMap  only just simulates how to deal with heatBeatCache without doing anything meaningful
func ProcessHeartBeatMap(heatBeatCache map[string]*data.HeartBeat) {
	heatBeatFilter := make(map[string]*controller.FilterInfo)
	observerdPods := make(map[string]map[string]*data.HeartBeatPod)
	nodeStatus := make(map[string]*v1alpha1.KoleQueryStatus)

	for i, hb := range heatBeatCache {
		heatBeatFilter[i] = &controller.FilterInfo{
			SeqNum:    hb.SeqNum,
			TimeStamp: hb.TimeStamp,
		}
		nodeStatus[hb.Name] = &v1alpha1.KoleQueryStatus{
			ObjectStatus: hb.State,
			ObjectName:   hb.Name,
			ObjectType:   v1alpha1.KoleObjectNode,
		}
		observerdPods[hb.Name] = make(map[string]*data.HeartBeatPod)
		for _, hbp := range hb.Pods {
			observerdPods[hb.Name][hbp.Key()] = &data.HeartBeatPod{
				Hash:      hbp.Hash,
				Name:      hbp.Name,
				NameSpace: hbp.NameSpace,
				Status:    hbp.Status,
			}
		}
	}
}
func (n *Compression) Run() error {
	//delete cr in DstNs namespace
	if err := n.DeleteSummariesByNS(n.DstNs); err != nil {
		return err
	}

	//list cr in SrcNs namespace to cache
	klog.V(4).Infof("Start to get source namespace cr to cache...")
	//flag default false,it means list cr in SrcNs namespace
	loadData, err := n.LoadSummariesByNs(n.SrcNs)
	if err != nil {
		klog.Errorf("Get source context err: ", err)
		return err
	}

	klog.Infof("#### Start to test ####")
	now := time.Now()
	totalNum, err := n.CreateSummariesByCache(n.DstNs, "nocompress", loadData)
	if err != nil {
		klog.Errorf("Create summaries to %s namespace fail: %v", n.DstNs, err)
		return err
	}
	needTime := time.Now().Sub(now).Milliseconds()
	klog.Infof("[No Compress] Create %d summaries to [%s] ns successfully, use %d ms", totalNum, n.DstNs, needTime)

	if n.OlSV {
		return nil
	}

	now = time.Now()
	heatBeatCache := make(map[string]*data.HeartBeat)
	cachedata, err := n.LoadSummariesByNs(n.DstNs)
	if err != nil {
		klog.Errorf("List cr in DstNs namespace fail: %v", err)
		return err
	}
	if err = json.Unmarshal(cachedata, &heatBeatCache); err != nil {
		klog.Errorf("unmarshal error %v", err)
		return err
	}
	ProcessHeartBeatMap(heatBeatCache)
	needTime = time.Now().Sub(now).Milliseconds()
	klog.Infof("[No Compress] Load all summaries from %s ns , datalen %d use %d ms", n.DstNs, len(cachedata), needTime)

	//delete cr in DstNs namespace
	if err := n.DeleteSummariesByNS(n.DstNs); err != nil {
		return err
	}

	//调用解压缩算法
	dataProcesser, err := options.AlgthmFactory(n.Algthm)
	if err != nil {
		return err
	}

	now = time.Now()
	cmp, err := dataProcesser.Compress(loadData)
	if err != nil {
		klog.Errorf("Compress err: %v", err)
		return err
	}
	var ratio = float64(len(loadData)-len(cmp)) / float64(len(loadData))

	totalNum, err = n.CreateSummariesByCache(n.DstNs, "compress", cmp)
	if err != nil {
		klog.Errorf("Create cr in DstNs namespace fail: %v", err)
		return err
	}
	needTime = time.Now().Sub(now).Milliseconds()
	klog.Infof("[Algthm %s][Compress] Create %d crs in %s namespace use %d ms. ratio %f origin data(len %d), commpressed data len %d",
		n.Algthm,
		totalNum,
		n.DstNs,
		needTime,
		ratio,
		len(loadData),
		len(cmp),
	)

	now = time.Now()
	res, err := n.LoadSummariesByNs(n.DstNs)
	if err != nil {
		klog.Errorf("List cr in DstNs namespace fail: %v", err)
		return err
	}

	uncmp, err := dataProcesser.UnCompress(res)
	if err != nil {
		klog.Errorf("Uncompress cache from cr in DstNs namespace fail: %v", err)
		return err
	}

	heatBeatCache = make(map[string]*data.HeartBeat)
	if err = json.Unmarshal(uncmp, &heatBeatCache); err != nil {
		klog.Errorf("unmarshal error %v", err)
		return err
	}
	ProcessHeartBeatMap(heatBeatCache)
	needTime = time.Now().Sub(now).Milliseconds()
	klog.Infof("[Algthm %s] Uncompress load cache from %s namespace use %d ms",
		n.Algthm,
		n.DstNs,
		needTime)

	//Finnal delete cr in DstNs namespace
	if err := n.DeleteSummariesByNS(n.DstNs); err != nil {
		return err
	}
	klog.Infof("Delete all summaries from %s namespace", n.DstNs)

	return nil
}

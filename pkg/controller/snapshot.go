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
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/kole-controller/app/options"
	"github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
	"github.com/openyurtio/kole/pkg/client/clientset/versioned"
	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/util"
)

func (c *KoleController) HeartBeatStatisticalLoop() {
	var interval time.Duration = 60 * 2
	time.Sleep(time.Second * interval)
	//wait.Forever(c.HeartBeatStatistical, time.Second*interval)
}

func (c *KoleController) SnapShotLoop() {

	ticker := time.NewTicker(time.Second * time.Duration(c.SnapshotInterval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.snapShot()
		}
	}
}

func (c *KoleController) GetReceiveNum() {
	klog.V(4).Infof("Get HB num %d", c.ReceiveNum)
}

func (c *KoleController) snapShot() {

	klog.Infof("Snapshot loop start ...")

	var registeringNum, registedNum, offlineNum int
	offlineMaps := make(map[int]int)
	nameToStatus := make(map[string]*v1alpha1.QueryNodeStatus)
	var hdata []byte
	var err error

	ackLists := make([]*data.HeartBeatACK, 0, 10000)

	n := time.Now().Unix()
	c.HeartBeatCache.SafeReadOperate(func() {
		if c.FirstSnapTime == 0 {
			c.FirstSnapTime = n
		}

		for _, hb := range c.HeartBeatCache.Cache {

			subTime := n - hb.LasterTimeStamp
			if hb.State == data.HeartBeatRegisterd && subTime >= c.HeartBeatTimeOut {
				klog.V(5).Infof("Nodename %s set offline, offline Time %d s", hb.Name, subTime)
				hb.State = data.HeartBeatOffline
			}

			if hb.State == data.HeartBeatRegistering {
				hb.State = data.HeartBeatRegisterd
				ackLists = append(ackLists, &data.HeartBeatACK{
					Identifier: hb.Identifier,
					Registerd:  true,
					NodeName:   hb.Name,
				})
				klog.V(5).Infof("Snapshot loop: find need ack hb[%s][%s]", hb.Identifier, hb.Name)
			}

			nameToStatus[hb.Name] = &v1alpha1.QueryNodeStatus{
				Status:          hb.State,
				InfEdgeNodeName: hb.Name,
			}

			i := 1
			for ; i <= 30; i++ {
				if subTime < int64(60*i) {
					break
				}
			}
			i = i - 1
			if _, ok := offlineMaps[i]; !ok {
				offlineMaps[i] = 1
			} else {
				offlineMaps[i] = offlineMaps[i] + 1
			}

			switch hb.State {
			case data.HeartBeatRegistering:
				registeringNum++
			case data.HeartBeatRegisterd:
				registedNum++
			case data.HeartBeatOffline:
				offlineNum++
			}
		}
		hdata, err = json.Marshal(c.HeartBeatCache.Cache)
		if err != nil {
			klog.Errorf("Snapshot Loop: marshal heartBeatCache error %v", hdata)
			return
		}
	})

	// Lock
	c.QueryNodeStatusCache.Reset(nameToStatus)

	if c.DataProcess != nil {
		hdata, err = c.DataProcess.Compress(hdata)
		if err != nil {
			klog.Info()
		}
	}
	c.syncAcks(ackLists)
	c.syncSummaris(hdata)

	var needTime int64
	nt := time.Now().Unix()
	if c.LasterSnapTime != 0 {
		needTime = nt - c.LasterSnapTime
	}
	klog.Infof("Snapshot Loop: registeringNum %d registerdNum %d offlineNum %d allNum %d len of HBCacheData is %d",
		registeringNum, registedNum, offlineNum, registedNum+registeringNum+offlineNum, len(hdata))
	klog.Infof("Current snap use %d s, laster jiange %d s, total jiange %d s", nt-n, needTime, nt-c.FirstSnapTime)

	for t, n := range offlineMaps {
		klog.Infof("Snapshot Loop: offline time %d m , node nums %d", t, n)
	}

	c.LasterSnapTime = nt
	c.LasterSnapIndex++
	klog.Infof("Snapshot Loop end ...")
}

func (c *KoleController) syncAcks(acks []*data.HeartBeatACK) {

	for i, _ := range acks {
		go func(ack *data.HeartBeatACK) {
			ctlTopic := filepath.Join(util.TopicCTLPrefix, ack.NodeName)
			if err := c.MessageHandler.PublishAck(context.Background(), ctlTopic, 0, false, ack); err != nil {
				klog.Errorf("Mqtt5 publish error %v", err)
				return
			}
		}(acks[i])

	}
}

func (c *KoleController) syncSummaris(hdata []byte) {
	snapedSummarisNames := make([]string, 0, 1024)
	namesLock := &sync.Mutex{}

	klog.V(4).Infof("Snapshot Loop: prepare to update summary ... ")
	deleteSummary := func(ns, name string) {
		for i := 0; i < 3; i++ {
			if err := c.LiteClient.LiteV1alpha1().Summaries(ns).Delete(context.Background(),
				name, metav1.DeleteOptions{}); err != nil {
				klog.Errorf("Delete[%d] old summary %s crd error %v", i, name, err)
				time.Sleep(time.Millisecond * 10)
			} else {
				break
			}
		}
	}
	deleteGroup := sync.WaitGroup{}
	for _, oldN := range c.SnapdSummaryNames {
		deleteGroup.Add(1)
		go func(ns, name string) {
			defer deleteGroup.Done()
			deleteSummary(ns, name)
		}(c.SummaryNS, oldN)
	}
	deleteGroup.Wait()

	// break down chunk
	bf := bytes.NewBuffer(hdata)
	lb := make(map[string]string)
	flag := fmt.Sprintf("%d", c.LasterSnapIndex)
	lb[util.SNAPSHOT_LABEL_IDENTIFIER] = flag
	lb[util.SNAPSHOT_LABEL_SUMMARY] = util.SNAPSHOT_LABEL_SUMMARY_VALUE
	bufferLen := util.SNAPSHOT_MAX_BUFFER_LEN

	if bf.Len()%bufferLen == 0 {
		lb[util.SNAPSHOT_LABEL_MAX_NUM] = fmt.Sprintf("%d", bf.Len()/bufferLen)
	} else {
		lb[util.SNAPSHOT_LABEL_MAX_NUM] = fmt.Sprintf("%d", bf.Len()/bufferLen+1)
	}

	createSummary := func(s *v1alpha1.Summary) {
		for j := 0; j < 3; j++ {
			if _, err := c.LiteClient.LiteV1alpha1().Summaries(s.Namespace).Create(context.Background(),
				s, metav1.CreateOptions{}); err != nil {
				klog.Errorf("create summary [%s][%s] error %v", s.GetNamespace(), s.GetName(), err)
				time.Sleep(time.Second)
			} else {
				klog.V(4).Infof("create summary [%s][%s] successful", s.GetNamespace(), s.GetName())
				namesLock.Lock()
				snapedSummarisNames = append(snapedSummarisNames, s.GetName())
				namesLock.Unlock()
				break
			}
		}

	}

	createGroup := sync.WaitGroup{}
	for i := 0; ; i++ {
		data := bf.Next(bufferLen)
		if len(data) == 0 {
			break
		}
		sum := &v1alpha1.Summary{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: c.SummaryNS,
				Name:      fmt.Sprintf("%s-%d", flag, i),
				Labels:    lb,
			},
			Data:  data,
			Index: i,
		}
		// TODO  need batch create
		createGroup.Add(1)
		go func(s *v1alpha1.Summary) {
			defer createGroup.Done()
			createSummary(s)
		}(sum)
	}
	createGroup.Wait()

	c.SnapdSummaryNames = snapedSummarisNames
	klog.Infof("Save %d summares successful", len(c.SnapdSummaryNames))
}

func LoadSnapShot(liteClient versioned.Interface, config *options.KoleControllerFlags, process DataProcesser) (
	map[string]*data.HeartBeat,
	map[string]*FilterInfo,
	[]string,
	map[string]map[string]*data.HeartBeatPod,
	map[string]*v1alpha1.QueryNodeStatus,
	error) {

	klog.Infof("Load snapshot start ...")

	heartBeatCache := make(map[string]*data.HeartBeat)
	heartBeatFilter := make(map[string]*FilterInfo)
	observerdPods := make(map[string]map[string]*data.HeartBeatPod)
	nodeStatus := make(map[string]*v1alpha1.QueryNodeStatus)
	// get current summery crd

	var timeoutS int64 = 60
	var continueStr string
	var max int64 = 500
	var total int

	hbData := make([]byte, 0, 1024*1000*1000)
	snapedName := make([]string, 0, 1024)
	allSummaris := make([]v1alpha1.Summary, 0, 1024)

	for {
		localSummaries, err := liteClient.LiteV1alpha1().Summaries(config.NameSpace).List(context.Background(), metav1.ListOptions{
			TimeoutSeconds: &timeoutS,
			Limit:          max,
			Continue:       continueStr,
		})
		if err != nil {
			klog.Errorf("List all summarys in ns[%s] error %v", config.NameSpace, err)
			return nil, nil, snapedName, observerdPods, nodeStatus, err
		}
		getLen := len(localSummaries.Items)
		total = total + getLen
		if getLen == 0 {
			break
		}

		// TODO 可以使用bytes.buffer
		for i, load := range localSummaries.Items {
			//data, ok := load.BinaryData[CONFIGMAP_KEY]
			snapedName = append(snapedName, load.GetName())
			allSummaris = append(allSummaris, localSummaries.Items[i])
		}

		//klog.Infof("Current get list node chunk nums %d ", len(listedNodes.Items))
		continueStr = localSummaries.GetContinue()

		if len(continueStr) == 0 {
			klog.Infof("continueStr is null , break ")
			break
		}

		//klog.Infof("Current get all node chunk nums %d,  continueStr %s ", allNum, continueStr)
		if getLen < int(max) {
			break
		}

	}

	// sort by index
	sort.Stable(BySummary(allSummaris))
	for i, _ := range allSummaris {
		hbData = append(hbData, allSummaris[i].Data...)
	}

	if total == 0 || len(hbData) == 0 {
		klog.Infof("Can not get any summary cr")
		return heartBeatCache, heartBeatFilter, snapedName, observerdPods, nodeStatus, nil
	}

	if process != nil {
		hbData, _ = process.UnCompress(hbData)
	}
	// TODO 可以使用fast json
	if err := json.Unmarshal(hbData, &heartBeatCache); err != nil {
		klog.Errorf("unmarshal error %v", err)
		return nil, nil, snapedName, nil, nil, err
	}
	for i, hb := range heartBeatCache {
		heartBeatFilter[i] = &FilterInfo{
			SeqNum:    hb.SeqNum,
			TimeStamp: hb.TimeStamp,
		}
		nodeStatus[hb.Name] = &v1alpha1.QueryNodeStatus{
			Status:          hb.State,
			InfEdgeNodeName: hb.Name,
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

	klog.Infof("Load snapshot end ...\n")
	return heartBeatCache, heartBeatFilter, snapedName, observerdPods, nodeStatus, nil
}

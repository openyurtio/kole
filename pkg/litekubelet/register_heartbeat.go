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
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/pkg/cache"
	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/util"
)

func (l *LiteKubelet) initHeartBeat() *data.HeartBeat {
	hb := util.InitMockHeartBeat(l.HostnameOverride, l.SeqNum)
	/*
		pods := make([]*data.HeartBeatPod, 0, 10)
		hash, err := kolecontroller.Md5PodSpec(&v1alpha1.PodSpec{
			Image:        "nginx.",
			Command:      []string{"/bin/bash", "start"},
			NodeSelector: map[string]string{"name": "node"},
		})
		if err != nil {
			return hb
		}
		for i := 0; i < 10; i++ {
			p := &data.HeartBeatPod{
				Hash:      hash,
				Name:      fmt.Sprintf("%s", uuid.New()),
				NameSpace: "infedge",
			}
			pods = append(pods, p)
		}
		hb.Pods = pods

	*/
	return hb
}

func (l *LiteKubelet) sendHeartBeat(hb *data.HeartBeat, qos byte, needACK bool) error {

	topic := util.TopicHeartBeat

	if err := l.MessageHandler.PublishData(context.Background(), topic, 0, false, hb); err != nil {
		return err
	}
	hbdata, err := json.Marshal(hb)
	if err != nil {
		klog.Errorf("Marshal Node to data error %v", err)
		return err
	}
	if needACK {
		klog.V(5).Infof("%s has send topic %s data , prepare to ack", hb.Identifier, topic)
		ack, ok := cache.GetDefaultTimeoutCache().PopWait(hb.Identifier, time.Second*time.Duration(l.HeartBeatInterval))
		if !ok {
			return fmt.Errorf("registering time out: node %s indentifier %s send topic %s, state %s", hb.Name, hb.Identifier, hb.State, topic)
		}
		if !ack.(*data.HeartBeatACK).Registerd {
			return fmt.Errorf("ack data is false: node %s indentifier %s send topic %s, state %s", hb.Name, hb.Identifier, hb.State, topic)
		}
		klog.V(4).Infof("#### data len %d , registering successful node %s", len(hbdata), hb.Name)
	} else {
		klog.V(5).Infof("@@@@ data len %d , registered successful node %s, len of pod %d", len(hbdata), hb.Name, len(hb.Pods))
	}
	return nil
}

func (l *LiteKubelet) registeringHeartBeat(needAck bool) (*data.HeartBeat, error) {
	hb := l.initHeartBeat()
	if err := l.sendHeartBeat(hb, 1, needAck); err != nil {
		return nil, err
	}
	return hb, nil
}

func (l *LiteKubelet) registerdHeartBeat(hb *data.HeartBeat) {
	hb.State = data.HeartBeatRegisterd
	hb.TimeStamp = time.Now().Unix()
	hb.Identifier = fmt.Sprintf("%v", uuid.New())
	l.sendHeartBeat(hb, 0, false)
}

func (l *LiteKubelet) syncHeartBeat(hb *data.HeartBeat) {
	localPodsLock.Lock()
	defer localPodsLock.Unlock()

	pods := make([]*data.HeartBeatPod, 0, 20)
	for _, p := range localPods {
		pp := &data.HeartBeatPod{
			Hash:      p.Hash,
			Name:      p.Name,
			NameSpace: p.NameSpace,
			Status: &data.HeartBeatPodStatus{
				Phase: data.HeartBeatPodStatusRunning,
			},
		}
		pods = append(pods, pp)
	}
	hb.Pods = pods
	return
}

func (l *LiteKubelet) registerdHeartBeatLoop(hb *data.HeartBeat) {
	ticker := time.NewTicker(time.Second * time.Duration(l.HeartBeatInterval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.syncHeartBeat(hb)
			l.registerdHeartBeat(hb)
		}
	}
}

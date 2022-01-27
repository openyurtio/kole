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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	outmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/lite-kubelet/app/options"
	"github.com/openyurtio/kole/pkg/data"
	"github.com/openyurtio/kole/pkg/message"
	"github.com/openyurtio/kole/pkg/util"
)

// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

type LiteKubelet struct {
	MessageHandler message.MessageHandler

	Registerd        bool
	IsMqtt5          bool
	MqttClient       outmqtt.Client
	MqttClient5      *autopaho.ConnectionManager
	Sub5CtlChan      chan *paho.Publish
	Sub5DataChan     chan *paho.Publish
	HostnameOverride string
	// s
	HeartBeatInterval int
	SubTopics         map[string]outmqtt.MessageHandler
	SeqNum            uint64
	IndexFlag         int
	ReceivePodDataNum int
}

type PersistentData struct {
	HostNameOverride string `json:"host_name_override"`
	SeqNum           uint64 `json:"seq_num"`
}

func NewMainLiteKubelet(deps *options.LiteKubeletFlags, index int, ismqtt5 bool) (*LiteKubelet, error) {
	hostname := os.Getenv("POD_NAME")
	if len(hostname) == 0 {
		hostname = fmt.Sprintf("%s", uuid.New())
	}
	hostnameOverride := fmt.Sprintf("%s-%d", hostname, index)

	seqNum, err := syncPersistentFile(deps.PersistentDir, hostnameOverride)
	if err != nil {
		return nil, err
	}

	lite := &LiteKubelet{
		HostnameOverride:  hostnameOverride,
		HeartBeatInterval: deps.HeartBeatInterval,
		SubTopics:         make(map[string]outmqtt.MessageHandler),
		SeqNum:            seqNum,
		Sub5CtlChan:       make(chan *paho.Publish, 1000),
		Sub5DataChan:      make(chan *paho.Publish, 1000),
		IsMqtt5:           ismqtt5,
		IndexFlag:         index,
	}

	if !deps.IsMqtt5 {
		// mqtt3
		h, err := message.NewMqtt3Handler(deps.Mqtt3Flags.MqttBroker, deps.Mqtt3Flags.MqttBrokerPort, deps.Mqtt3Flags.MqttInstance, deps.Mqtt3Flags.MqttGroup,
			hostnameOverride,
			map[string]outmqtt.MessageHandler{
				filepath.Join(util.TopicCTLPrefix, lite.HostnameOverride):  lite.SubCTL,
				filepath.Join(util.TopicDataPrefix, lite.HostnameOverride): lite.SubData,
			})
		if err != nil {
			return nil, err
		}
		lite.MessageHandler = h

	} else {
		// mqtt5
		// mqtt 5
		h, err := message.NewMqtt5Handler(deps.Mqtt5Flags.MqttServer, lite.CreateSubscribes5(), hostnameOverride, true)
		if err != nil {
			return nil, err
		}
		lite.MessageHandler = h
	}

	// 此调试信息不要修改， 跟shell 关联
	klog.V(4).Infof("--- Create client and Subscribe all topic successful node %s ---", hostnameOverride)

	return lite, nil
}

func (l *LiteKubelet) runRealyLoop() {
	var hb *data.HeatBeat
	var err error

	for {
		hb, err = l.registeringHeatBeat(true)
		if err != nil {
			klog.Errorf("%s Registering heatbeat error %v", l.HostnameOverride, err)
			time.Sleep(time.Second * 10)
			continue
		}
		break
	}
	l.Registerd = true

	l.registerdHeatBeatLoop(hb)
}

func (l *LiteKubelet) Run() {

	if l.IsMqtt5 {
		l.ConsumeSubLoop()
	}

	l.runRealyLoop()
}

func syncPersistentFile(dir, hostnameOverride string) (uint64, error) {

	var SeqNum uint64
	var data []byte
	var err error
	firstCreate := false

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		klog.Errorf("mkdir %s error %v", dir, err)
		return 0, err
	}

	file := filepath.Join(dir, fmt.Sprintf("%s.yaml", hostnameOverride))

	_, err = os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		firstCreate = true
		_, err := os.Create(file)
		if err != nil {
			klog.Errorf("create file %s error %v", file, err)
			return 0, err
		}
	}

	// must use O_RDWR, CREATE AND TRUNC
	if !firstCreate {
		data, err = ioutil.ReadFile(file)
		if err != nil {
			klog.Errorf("Read file %s data error %v", file, err)
			return 0, err
		}

		tmp := &PersistentData{}
		if err := json.Unmarshal(data, tmp); err != nil {
			klog.Errorf("Unarshal error %v", err)
			return 0, err
		}
		tmp.SeqNum = tmp.SeqNum + 1

		data, err = json.Marshal(tmp)
		if err != nil {
			klog.Errorf("Marshal added persistentData error %v", err)
			return 0, err
		}
		SeqNum = tmp.SeqNum
		klog.V(5).Infof("Write to file %s data:\n%s", file, string(data))
	} else {
		tmp := &PersistentData{
			HostNameOverride: hostnameOverride,
			SeqNum:           1,
		}
		data, err = json.Marshal(tmp)
		if err != nil {
			klog.Errorf("Marshal init persistentData error %v", err)
			return 0, err
		}
		SeqNum = tmp.SeqNum
		klog.V(5).Infof("First write to file %s data:\n%s", file, string(data))
	}

	if err := ioutil.WriteFile(file, data, 0644); err != nil {
		klog.Errorf("Flush data to file %s error %v", file, err)
		return 0, err
	}
	klog.V(5).Infof("Write to file %s success. data:\n%s", file, string(data))
	return SeqNum, nil
}

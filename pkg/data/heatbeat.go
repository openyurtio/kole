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
package data

import (
	"encoding/json"
	"fmt"
)

const HeatBeatRegistering = "Registering"
const HeatBeatRegisterd = "Registerd"
const HeatBeatOffline = "Offline"

const OFFLINE_TIMEOUT = 60 * 20

type HeatBeat struct {
	Name string `json:"name,omitempty"`
	// 未来可能用于workload 的label selector
	Labels    map[string]string `json:"labels,omitempty"`
	TimeStamp int64             `json:"timestamp,omitempty"`
	// 仅用于云端控制器记录上一次收到的时间戳记录，用于判断是否超时
	LasterTimeStamp int64           `json:"-"`
	Identifier      string          `json:"identifier,omitempty"`
	SeqNum          uint64          `json:"seqnum,omitempty"`
	State           string          `json:"state,omitempty"`
	Status          *HeatBeatStatus `json:"status,omitempty"`
	Pods            []*HeatBeatPod  `json:"pods,omitempty"`
}

type HeatBeatStatus struct {
	// IP地址
	Addresses []*Address `json:"addresses,omitempty"`
	// 可分配的资源， 便于云端控制器调度
	Allocatable *Resource `json:"allocatable,omitempty"`
	// 容量
	Capacity *Resource `json:"capacity,omitempty"`
	NodeInfo *NodeInfo `json:"nodeInfo,omitempty"`
}

const AddressTypeInternal = "InternalIP"
const AddressTypeHostName = "HostName"

type Address struct {
	Address string
	Type    string
}

type Resource struct {
	// m
	Cpu int `json:"cpu,omitempty"`
	// Ki
	Memory int `json:"memory,omitempty"`
	// pods num
	Pods int `json:"pods,omitempty"`
}

type NodeInfo struct {
	Architecture       string
	LiteKubeletVersion string
	KernelVersion      string
}

type HeatBeatPod struct {
	// 只需要上报它的hash 值， 这样方便云端控制器 根据workload 判断是否要重新下发新的配置
	// 不需要上传pod 的具体spec 信息
	Hash      string             `json:"hash,omitempty"`
	Name      string             `json:"name,omitempty"`
	NameSpace string             `json:"namespace,omitempty"`
	Status    *HeatBeatPodStatus `json:"status,omitempty"`
}

func (p *HeatBeatPod) Key() string {
	return fmt.Sprintf("%s-%s", p.NameSpace, p.Name)
}

const HeatBeatPodStatusRunning = "Running"

type HeatBeatPodStatus struct {
	// Runing ,Completed, Termaled ...
	Phase string `json:"phase,omitempty"`
}

func UnmarshalPayloadToHeatBeat(payload []byte) (*HeatBeat, error) {
	d := &HeatBeat{}
	if err := json.Unmarshal(payload, d); err != nil {
		return nil, err
	}
	return d, nil
}

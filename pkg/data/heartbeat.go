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

const HeartBeatRegistering = "Registering"
const HeartBeatRegisterd = "Registerd"
const HeartBeatOffline = "Offline"

const OFFLINE_TIMEOUT = 60 * 20

type HeartBeat struct {
	Name      string            `json:"name,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	TimeStamp int64             `json:"timestamp,omitempty"`
	// Used in cloud state cache only
	LasterTimeStamp int64            `json:"-"`
	Identifier      string           `json:"identifier,omitempty"`
	SeqNum          uint64           `json:"seqnum,omitempty"`
	State           string           `json:"state,omitempty"`
	Status          *HeartBeatStatus `json:"status,omitempty"`
	Pods            []*HeartBeatPod  `json:"pods,omitempty"`
}

type HeartBeatStatus struct {
	// IP
	Addresses []*Address `json:"addresses,omitempty"`
	// Resources that can be used by cloud scheduler
	Allocatable *Resource `json:"allocatable,omitempty"`
	Capacity    *Resource `json:"capacity,omitempty"`
	NodeInfo    *NodeInfo `json:"nodeInfo,omitempty"`
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

type HeartBeatPod struct {
	// We only report the hash of the Pod Spec in the edge node to compare against the spec in the cloud state cache.
	Hash      string              `json:"hash,omitempty"`
	Name      string              `json:"name,omitempty"`
	NameSpace string              `json:"namespace,omitempty"`
	Status    *HeartBeatPodStatus `json:"status,omitempty"`
}

func (p *HeartBeatPod) Key() string {
	return fmt.Sprintf("%s-%s", p.NameSpace, p.Name)
}

const HeartBeatPodStatusRunning = "Running"

type HeartBeatPodStatus struct {
	// Runing ,Completed, Termaled ...
	Phase string `json:"phase,omitempty"`
}

func UnmarshalPayloadToHeartBeat(payload []byte) (*HeartBeat, error) {
	d := &HeartBeat{}
	if err := json.Unmarshal(payload, d); err != nil {
		return nil, err
	}
	return d, nil
}

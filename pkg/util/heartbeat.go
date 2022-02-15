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
package util

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/openyurtio/kole/pkg/data"
)

func InitMockHeartBeat(hostname string, seqNum uint64) *data.HeartBeat {
	labels := make(map[string]string)
	labels["HostName"] = hostname

	return &data.HeartBeat{
		SeqNum:     seqNum,
		TimeStamp:  time.Now().Unix(),
		Identifier: fmt.Sprintf("%v", uuid.New()),
		State:      data.HeartBeatRegistering,
		Name:       hostname,
		Labels:     labels,
		Status: &data.HeartBeatStatus{
			Addresses: []*data.Address{
				{
					Address: "127.0.0.1",
					Type:    data.AddressTypeInternal,
				},
				{
					Address: "localhost",
					Type:    data.AddressTypeHostName,
				},
			},
			Allocatable: &data.Resource{
				Cpu:    6000,
				Memory: 6000 * 1024,
				Pods:   8,
			},
			Capacity: &data.Resource{
				Cpu:    8000,
				Memory: 8000 * 1024,
				Pods:   10,
			},
			NodeInfo: &data.NodeInfo{
				Architecture:       "amd64",
				LiteKubeletVersion: "v0.1.0",
				KernelVersion:      "4.19.91",
			},
		},
	}
}

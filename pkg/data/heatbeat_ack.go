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
)

type HeatBeatACK struct {
	Identifier string `json:"identifier,omitempty"`
	Registerd  bool   `json:"registerd,omitempty"`
	NodeName   string `json:"-"`
}

func UnmarshalPayloadToHeatBeatACK(payload []byte) (*HeatBeatACK, error) {
	d := &HeatBeatACK{}
	if err := json.Unmarshal(payload, d); err != nil {
		return nil, err
	}
	return d, nil
}

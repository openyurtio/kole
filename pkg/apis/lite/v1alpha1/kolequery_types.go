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
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,path=kolequeries,shortName=kq,categories=all
// +kubebuilder:subresource:status

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// KoleQuery is
type KoleQuery struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec *KoleQuerySpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// +optional
	Status []*KoleQueryStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

type KoleQueryType string

const (
	KoleQueryGet   KoleQueryType = "Get"
	KoleQueryWatch KoleQueryType = "Watch"
)

type KoleQueryObjectType string

const (
	KoleObjectNode KoleQueryObjectType = "Node"
	KoleObjectPod  KoleQueryObjectType = "Pod"
)

type KoleQuerySpec struct {
	// For now, the KoleQueryWatch type is only supported for querying a single object.
	// +optional
	QueryType KoleQueryType `json:"queryType,omitempty"`
	// +optional
	ObjectType KoleQueryObjectType `json:"objectType"`
	// +optional
	ObjectName string `json:"objectName,omitempty"`
	// +optional
	ObjectStatus string `json:"objectStatus,omitempty"`
	// label selectors of the query objects
	// +optional
	ObjectSelector map[string]string `json:"objectSelector,omitempty"`
}

type KoleQueryStatus struct {
	LastObservedTime metav1.Time         `json:"lastObservedTime"`
	ObjectType       KoleQueryObjectType `json:"objectType"`
	ObjectStatus     string              `json:"objectStatus"`
	ObjectName       string              `json:"objectName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QueryNodeList is
type KoleQueryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []KoleQuery `json:"items"`
}

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
// +kubebuilder:resource:scope=Namespaced,path=querynodes,shortName=qn,categories=all
// +kubebuilder:subresource:status

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// QueryNode is
type QueryNode struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec *QueryNodeSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// +optional
	Status []*QueryNodeStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

type QueryNodeSpec struct {
	// NodePoolSelector is a label query over nodepool that should match the replica count.
	// It must match the nodepool's labels.
	// +optional
	NodeName string `json:"nodeName,omitempty"`
	// +optional
	NodeStatus string `json:"nodeStatus,omitempty"`
	// +optional
	NodeLabelSelector string `json:"nodeLabelSelector,omitempty"`
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

type QueryNodeStatus struct {
	Status          string `json:"status"`
	InfEdgeNodeName string `json:"infEdgeNodeName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QueryNodeList is
type QueryNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []QueryNode `json:"items"`
}

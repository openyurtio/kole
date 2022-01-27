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
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"

	"github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
)

var infEdgeNodeTemplate, smallInfEdgeNodeTemplate *v1alpha1.InfEdgeNode
var InfEdgeNodeTemplateSize, SmallInfEdgeNodeTemplateSize int

func init() {

	infEdgeNodeTemplate = &v1alpha1.InfEdgeNode{}
	smallInfEdgeNodeTemplate = &v1alpha1.InfEdgeNode{}
	nodeStr := `
apiVersion: lite.openyurt.io/v1alpha1 
kind: InfEdgeNode 
metadata:
  annotations:
    csi.volume.kubernetes.io/nodeid: '{"diskplugin.csi.alibabacloud.com":"i-8vb2y0cdlriptkfer0ir","nasplugin.csi.alibabacloud.com":"i-8vb2y0cdlriptkfer0ir","ossplugin.csi.alibabacloud.com":"i-8vb2y0cdlriptkfer0ir"}'
    flannel.alpha.coreos.com/backend-data: "null"
    flannel.alpha.coreos.com/backend-type: ""
    flannel.alpha.coreos.com/kube-subnet-manager: "true"
    flannel.alpha.coreos.com/public-ip: 192.168.0.81
    kubeadm.alpha.kubernetes.io/cri-socket: /run/containerd/containerd.sock
    node.alpha.kubernetes.io/ttl: "0"
    volumes.kubernetes.io/controller-managed-attach-detach: "true"
  creationTimestamp: "2021-12-27T07:34:53Z"
  labels:
    ack.aliyun.com: cfcf334401e6845e69400792e4fb9128c
    alibabacloud.com/nodepool-id: np66c94fbe4b4d46f18fefc95f51757df7
    beta.kubernetes.io/arch: amd64
    beta.kubernetes.io/instance-type: ecs.g5.xlarge
    beta.kubernetes.io/os: linux
    failure-domain.beta.kubernetes.io/region: cn-zhangjiakou
    failure-domain.beta.kubernetes.io/zone: cn-zhangjiakou-c
    kubernetes.io/arch: amd64
    kubernetes.io/hostname: cn-zhangjiakou.192.168.0.81
    kubernetes.io/os: linux
    node.csi.alibabacloud.com/disktype.cloud_efficiency: available
    node.csi.alibabacloud.com/disktype.cloud_essd: available
    node.csi.alibabacloud.com/disktype.cloud_ssd: available
    node.kubernetes.io/instance-type: ecs.g5.xlarge
    nodepool: zhangjie
    topology.diskplugin.csi.alibabacloud.com/zone: cn-zhangjiakou-c
    topology.kubernetes.io/region: cn-zhangjiakou
    topology.kubernetes.io/zone: cn-zhangjiakou-c
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:labels:
          f:beta.kubernetes.io/instance-type: {}
          f:failure-domain.beta.kubernetes.io/region: {}
          f:failure-domain.beta.kubernetes.io/zone: {}
          f:node.kubernetes.io/instance-type: {}
          f:topology.kubernetes.io/region: {}
          f:topology.kubernetes.io/zone: {}
      f:status:
        f:conditions:
          k:{"type":"NetworkUnavailable"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
    manager: cloud-controller-manager
    operation: Update
    time: "2021-12-27T07:34:53Z"
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          f:node.alpha.kubernetes.io/ttl: {}
      f:spec:
        f:podCIDR: {}
        f:podCIDRs:
          .: {}
          v:"10.155.2.192/26": {}
    manager: kube-controller-manager
    operation: Update
    time: "2021-12-27T07:34:53Z"
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          f:kubeadm.alpha.kubernetes.io/cri-socket: {}
    manager: kubeadm
    operation: Update
    time: "2021-12-27T07:34:53Z"
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          f:flannel.alpha.coreos.com/backend-data: {}
          f:flannel.alpha.coreos.com/backend-type: {}
          f:flannel.alpha.coreos.com/kube-subnet-manager: {}
          f:flannel.alpha.coreos.com/public-ip: {}
    manager: flanneld
    operation: Update
    time: "2021-12-27T07:35:22Z"
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:status:
        f:conditions:
          k:{"type":"DockerOffline"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
          k:{"type":"InodesPressure"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
          k:{"type":"InstanceExpired"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
          k:{"type":"KernelDeadlock"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
          k:{"type":"NTPProblem"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
          k:{"type":"NodePIDPressure"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
          k:{"type":"ReadonlyFilesystem"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
          k:{"type":"RuntimeOffline"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
    manager: node-problem-detector
    operation: Update
    time: "2021-12-27T07:35:35Z"
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:labels:
          f:node.csi.alibabacloud.com/disktype.cloud_efficiency: {}
          f:node.csi.alibabacloud.com/disktype.cloud_essd: {}
          f:node.csi.alibabacloud.com/disktype.cloud_ssd: {}
    manager: plugin.csi.alibabacloud.com
    operation: Update
    time: "2021-12-27T07:35:39Z"
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .: {}
          f:csi.volume.kubernetes.io/nodeid: {}
          f:volumes.kubernetes.io/controller-managed-attach-detach: {}
        f:labels:
          .: {}
          f:ack.aliyun.com: {}
          f:alibabacloud.com/nodepool-id: {}
          f:beta.kubernetes.io/arch: {}
          f:beta.kubernetes.io/os: {}
          f:kubernetes.io/arch: {}
          f:kubernetes.io/hostname: {}
          f:kubernetes.io/os: {}
          f:nodepool: {}
          f:topology.diskplugin.csi.alibabacloud.com/zone: {}
      f:spec:
        f:providerID: {}
      f:status:
        f:addresses:
          .: {}
          k:{"type":"Hostname"}:
            .: {}
            f:address: {}
            f:type: {}
          k:{"type":"InternalIP"}:
            .: {}
            f:address: {}
            f:type: {}
        f:allocatable:
          .: {}
          f:cpu: {}
          f:ephemeral-storage: {}
          f:hugepages-1Gi: {}
          f:hugepages-2Mi: {}
          f:memory: {}
          f:pods: {}
        f:capacity:
          .: {}
          f:cpu: {}
          f:ephemeral-storage: {}
          f:hugepages-1Gi: {}
          f:hugepages-2Mi: {}
          f:memory: {}
          f:pods: {}
        f:conditions:
          .: {}
          k:{"type":"DiskPressure"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
          k:{"type":"MemoryPressure"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
          k:{"type":"PIDPressure"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
          k:{"type":"Ready"}:
            .: {}
            f:lastHeartbeatTime: {}
            f:lastTransitionTime: {}
            f:message: {}
            f:reason: {}
            f:status: {}
            f:type: {}
        f:config: {}
        f:daemonEndpoints:
          f:kubeletEndpoint:
            f:Port: {}
        f:images: {}
        f:nodeInfo:
          f:architecture: {}
          f:bootID: {}
          f:containerRuntimeVersion: {}
          f:kernelVersion: {}
          f:kubeProxyVersion: {}
          f:kubeletVersion: {}
          f:machineID: {}
          f:operatingSystem: {}
          f:osImage: {}
          f:systemUUID: {}
    manager: kubelet
    operation: Update
    time: "2021-12-27T07:35:40Z"
  name: cn-zhangjiakou.192.168.0.81
  resourceVersion: "2768192"
  uid: ba1089c4-282c-47b4-99c6-d8810361de5a
spec:
  podCIDR: 10.155.2.192/26
  podCIDRs:
  - 10.155.2.192/26
  providerID: cn-zhangjiakou.i-8vb2y0cdlriptkfer0ir
status:
  addresses:
  - address: 192.168.0.81
    type: InternalIP
  - address: cn-zhangjiakou.192.168.0.81
    type: Hostname
  allocatable:
    cpu: 3900m
    ephemeral-storage: "37926431477"
    hugepages-1Gi: "0"
    hugepages-2Mi: "0"
    memory: 13152392Ki
    pods: "64"
  capacity:
    cpu: "4"
    ephemeral-storage: 41152812Ki
    hugepages-1Gi: "0"
    hugepages-2Mi: "0"
    memory: 16117896Ki
    pods: "64"
  conditions:
  - lastHeartbeatTime: "2021-12-30T09:24:36Z"
    lastTransitionTime: "2021-12-27T07:35:34Z"
    message: ntp service is up
    reason: NTPIsUp
    status: "False"
    type: NTPProblem
  - lastHeartbeatTime: "2021-12-30T09:24:36Z"
    lastTransitionTime: "2021-12-27T07:35:34Z"
    message: node has no inodes pressure
    reason: NodeHasNoInodesPressure
    status: "False"
    type: InodesPressure
  - lastHeartbeatTime: "2021-12-30T09:24:36Z"
    lastTransitionTime: "2021-12-27T07:35:34Z"
    message: instance is not going to be terminated
    reason: InstanceNotToBeTerminated
    status: "False"
    type: InstanceExpired
  - lastHeartbeatTime: "2021-12-30T09:24:36Z"
    lastTransitionTime: "2021-12-27T07:35:34Z"
    message: kernel has no deadlock
    reason: KernelHasNoDeadlock
    status: "False"
    type: KernelDeadlock
  - lastHeartbeatTime: "2021-12-30T09:24:36Z"
    lastTransitionTime: "2021-12-27T07:35:34Z"
    message: Filesystem is read-only
    reason: FilesystemIsReadOnly
    status: "False"
    type: ReadonlyFilesystem
  - lastHeartbeatTime: "2021-12-30T09:24:36Z"
    lastTransitionTime: "2021-12-27T07:35:34Z"
    message: Node has no PID Pressure
    reason: NodeHasNoPIDPressure
    status: "False"
    type: NodePIDPressure
  - lastHeartbeatTime: "2021-12-30T09:24:36Z"
    lastTransitionTime: "2021-12-27T07:35:34Z"
    message: docker daemon is ok
    reason: DockerDaemonNotOffline
    status: "False"
    type: DockerOffline
  - lastHeartbeatTime: "2021-12-30T09:24:36Z"
    lastTransitionTime: "2021-12-27T07:35:34Z"
    message: container runtime daemon is ok
    reason: RuntimeDaemonNotOffline
    status: "False"
    type: RuntimeOffline
  - lastHeartbeatTime: "2021-12-30T09:27:52Z"
    lastTransitionTime: "2021-12-27T07:34:53Z"
    message: kubelet has sufficient memory available
    reason: KubeletHasSufficientMemory
    status: "False"
    type: MemoryPressure
  - lastHeartbeatTime: "2021-12-30T09:27:52Z"
    lastTransitionTime: "2021-12-27T07:34:53Z"
    message: kubelet has no disk pressure
    reason: KubeletHasNoDiskPressure
    status: "False"
    type: DiskPressure
  - lastHeartbeatTime: "2021-12-30T09:27:52Z"
    lastTransitionTime: "2021-12-27T07:34:53Z"
    message: kubelet has sufficient PID available
    reason: KubeletHasSufficientPID
    status: "False"
    type: PIDPressure
  - lastHeartbeatTime: "2021-12-30T09:27:52Z"
    lastTransitionTime: "2021-12-27T07:35:25Z"
    message: kubelet is posting ready status
    reason: KubeletReady
    status: "True"
    type: Ready
  - lastHeartbeatTime: "2021-12-27T07:34:53Z"
    lastTransitionTime: "2021-12-27T07:34:53Z"
    message: RouteController created a route
    reason: RouteCreated
    status: "False"
    type: NetworkUnavailable
  config: {}
  daemonEndpoints:
    kubeletEndpoint:
      Port: 10250
  images:
  - names:
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/csi-plugin@sha256:1bdc9a7d06451538fd409069982bcd03fbd492ab06fde0673fad7c5ad6589c14
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/csi-plugin:v1.20.7-aafce42-aliyun
    sizeBytes: 251614931
  - names:
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/logtail@sha256:da549e535549b9124f072cca0e0984643dda4d58234671b9250b368708ead153
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/logtail:v0.16.62.3-da583e0-aliyun
    sizeBytes: 158136049
  - names:
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/node-problem-detector@sha256:44fbc8400b7ed267aef765714e3423a66721fc04585ee701c0e9e5182cec796e
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/node-problem-detector:v0.8.10-e0ff7d2
    sizeBytes: 115728383
  - names:
    - registry.cn-beijing.aliyuncs.com/infedge/lite-kubelet@sha256:b38deb16c36d559d0ee6b64ce10cf49bb05f1452a02fc1143d792254f95a37af
    - registry.cn-beijing.aliyuncs.com/infedge/lite-kubelet:v1
    sizeBytes: 81017449
  - names:
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/kube-proxy@sha256:51dbd5ab4a0aeaf62970f889a930243e44eedbf3cefc825120486344d8404cf9
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/kube-proxy:v1.20.11-aliyun.1
    sizeBytes: 46089579
  - names:
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/k8s-dns-node-cache@sha256:2004485b75f2900f7e6ca18de7780e968ac2267442a9bdcc1329675ed2592732
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/k8s-dns-node-cache:v1.15.13-6-7e6778ac
    sizeBytes: 42154699
  - names:
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/flannel@sha256:6911127cf398562cff7aedec874c3e5fb76302b0572cb0dc187cd0d3ba626c30
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/flannel:v0.13.0.2-466064b-aliyun
    sizeBytes: 21951119
  - names:
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/kube-rbac-proxy@sha256:d2872c6a945494e5929de3785eb0c9232a9868ec1dc8e45a58bb08b815ec09ad
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/kube-rbac-proxy:v0.4.1-bbc79f2e-aliyun
    sizeBytes: 17155877
  - names:
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/csi-node-driver-registrar@sha256:e3a4959ba4843e3c470d077edc759324078e6781289c4a622a80d7974741eee1
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/csi-node-driver-registrar:v1.3.0-6e9fff3-aliyun
    sizeBytes: 7716815
  - names:
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/arms-docker-repo/node-exporter@sha256:ec3c058e92915388f7abb5d767e3d3f3712f21214c90dfa01fd42c43feb3c3d1
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/arms-docker-repo/node-exporter:v0.17.0-slim
    sizeBytes: 7147075
  - names:
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/busybox@sha256:3058e3a1129c64da64d5c7889e6eedb0666262d7ee69b289f2d4379f69362383
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/busybox:v1.29.2
    sizeBytes: 712890
  - names:
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/pause@sha256:1ff6c18fbef2045af6b9c16bf034cc421a29027b800e4f9b68ae9b1cb3e9ae07
    - registry-vpc.cn-zhangjiakou.aliyuncs.com/acs/pause:3.5
    sizeBytes: 301416
  nodeInfo:
    architecture: amd64
    bootID: 49020e5f-ac30-425c-833f-819222e90915
    containerRuntimeVersion: containerd://1.4.8
    kernelVersion: 4.19.91-24.1.al7.x86_64
    kubeProxyVersion: v1.20.11-aliyun.1
    kubeletVersion: v1.20.11-aliyun.1
    machineID: "20210726155106994708212788883726"
    operatingSystem: linux
    osImage: Alibaba Cloud Linux (Aliyun Linux) 2.1903 LTS (Hunting Beagle)
    systemUUID: fa86a3b8-31a7-45d6-b4ae-4bf17b42800b
`
	if err := yaml.Unmarshal([]byte(nodeStr), infEdgeNodeTemplate); err != nil {
		klog.Fatalf("Unmarshal error %v", err)
	}
	data, err := yaml.Marshal(infEdgeNodeTemplate)
	if err != nil {
		klog.Fatalf("Marshal infEdgeNodeTemplate to data error %v", err)
	}
	InfEdgeNodeTemplateSize = len(data) / 1024

	var smallNodeStr = `
apiVersion: v1
kind: Node
metadata:
  annotations:
    flannel/public-ip: 192.168.0.81
  labels:
    node.csi/cloud_ssd: available
  name: cn-zhangjiakou.192.168.0.81
  resourceVersion: "2768192"
  uid: ba1089c4-282c-47b4-99c6-d8810361de5a
spec:
  podCIDR: 10.155.2.192/26
status:
  allocatable:
    cpu: 3900m
    memory: 13152392Ki
  capacity:
    cpu: "4"
    memory: 16117896Ki
  conditions:
  - lastHeartbeatTime: "2021-12-30T09:24:36Z"
  daemonEndpoints:
    kubeletEndpoint:
      Port: 10250
  nodeInfo:
    kernelVersion: 4.19.91-24.1.al7.x86_64
    kubeProxyVersion: v1.20.11-aliyun.1
    kubeletVersion: v1.20.11-aliyun.1
    machineID: "20210726155106994708212788883726"
    operatingSystem: linux
`

	if err := yaml.Unmarshal([]byte(smallNodeStr), smallInfEdgeNodeTemplate); err != nil {
		klog.Fatalf("Unmarshal error %v", err)
	}
	smallData, err := yaml.Marshal(smallInfEdgeNodeTemplate)
	if err != nil {
		klog.Fatalf("Marshal smallInfEdgeNodeTemplate to data error %v", err)
	}
	SmallInfEdgeNodeTemplateSize = len(smallData) / 1024
}

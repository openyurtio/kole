# Instructions

This tutorial focuses on how to deploy the Kole project

## Environmental requirements

### Kubernetes Cluster

Need to create a kubernetes ahead of the cluster, here we use the ali cloud container service, https://cs.console.aliyun.com/, provides a standard  Pro cluster.

### Emqx 集群
Emqx clusters can be deployed on kubernetes clusters. For details about the deployment mode, see [link](https://www.emqx.com/zh/blog/rapidly-deploy-emqx-clusters-on-kubernetes-via-helm)

It is recommended that the emqx service type be LoadBalancer. The LoadBalancer service provides a public IP address (LOADBALANCE_IP). Next we will use the LOADBALANCE_IP IP。

# How to deploy

## Deploy Instructions

All of the following operations are in the kole source root directory.

##  Deploy Kole related components

Kole consists of two main components, kole-controller and lite-kubelet, as well as some CRD definitions.All related components are deployed under the namespace kole.

``` 
$ kubectl apply -f config/setup/manifest.yaml

namespace/kole created
service/kole-controller created
statefulset.apps/kole-controller created
configmap/lite-kubelet-start-signal created
service/lite-kubelet created
statefulset.apps/lite-kubelet created
customresourcedefinition.apiextensions.k8s.io/infdaemonsets.lite.openyurt.io created
customresourcedefinition.apiextensions.k8s.io/infedgenodes.lite.openyurt.io created
customresourcedefinition.apiextensions.k8s.io/querynodes.lite.openyurt.io created
customresourcedefinition.apiextensions.k8s.io/summaries.lite.openyurt.io created
clusterrole.rbac.authorization.k8s.io/kole created
clusterrolebinding.rbac.authorization.k8s.io/kole-rolebinding created
```

Check out the statefulsets of Kole-Controller and Lite-Kubelet:

``` 
$ kubectl get statefulsets.apps -n kole
NAME              READY   AGE
kole-controller   0/0     92s
lite-kubelet      0/0     92s
```

Change the value of the statefulSet environment variable MQTT5_SERVER for kole-controller and lite-kubelet.

```
        env:
        - name: MQTT5_SERVER
          value: mqtt://8.142.157.229:1883
```

We need to put the MQTT: / / 8.142.157.229:1883 this value into the new emqx service corresponding to the public IP address: MQTT: / / ${LOADBALANCE_IP} :1883

Run the following command to change the value:

```
CONTROLLER_PATCH=$(cat <<- EOF
spec:
  template:
    spec:
      containers:
      - name: kole-controller
        env:
        - name: MQTT5_SERVER
          value: "mqtt://{LOADBALANCE_IP}:1883"
EOF
)

 $ kubectl patch -n kole statefulsets.apps kole-controller  --patch "$CONTROLLER_PATCH"
 
 LITE_KUBELET_PATCH=$(cat <<- EOF
spec:
  template:
    spec:
      containers:
      - name: lite-kubelet
        env:
        - name: MQTT5_SERVER
          value: "mqtt://{LOADBALANCE_IP}:1883"
EOF
)
 
$ kubectl patch -n kole statefulsets.apps lite-kubelet --patch "$CONTROLLER_PATCH"

```

After the modification is complete, change the number of kole-controller instances to 1 and the number of lite-kubelet instances to 1. The number of lite-kubelet instances can be multiple.

```
 $ kubectl scale statefulsets.apps -n kole kole-controller --replicas 1
 
 $ kubectl scale statefulsets.apps -n kole lite-kubelet --replicas 1
```

You can use Kubectl logs to view kole-Controller and Lite-kubelet logs:

```
$ kubectl -n kole logs -f kole-controller-0

$ kubectl -n kole logs -f lite-kubelete-0
```

In the Kole namespace, we create a lite-kubelet-start-signal configmap resource. If we change any configuration of this Configmap, All lite-kubelet MQTT connections are triggered for registring HB and registerd HB operations.

```
$ kubectl get cm -n kole  lite-kubelet-start-signal -o yaml
apiVersion: v1

kind: ConfigMap
metadata:
  creationTimestamp: "2022-01-27T08:29:15Z"
  name: lite-kubelet-start-signal
  namespace: kole
  resourceVersion: "4193091"
  uid: 5b9f614c-74c0-4679-8d3d-458c245b587f
  
data:
  test: start
```

We can modify the configuration using the following command：

```
CM_PATCH=$(cat <<- EOF
data:
  test: start1111
EOF
)

$ kubectl patch -n kole cm lite-kubelet-start-signal  --patch "$CM_PATCH"
```

Check the log of kole-Controller to see the number of registered hosts：

```
$ kubectl -n kole logs -f kole-controller-0
```







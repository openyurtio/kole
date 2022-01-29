#!/bin/bash

set +x

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
YAML_FILE="$KOLE_OUTPUT_DIR/manifest.yaml"

CONFIG_BASE_DIR="$KOLE_ROOT/config/base"

CONFIG_CRD_DIR="$KOLE_ROOT/config/crd"
CONFIG_RBAC_DIR="$KOLE_ROOT/config/rbac"

function create_basefile(){

cat << EOF >$YAML_FILE
apiVersion: v1
kind: Namespace
metadata:
  name: kole
---

apiVersion: v1
kind: Service
metadata:
  name: kole-controller 
  namespace: kole 
  labels:
    app: kole-controller 
spec:
  ports:
  - port: 80
    name: kole-controller 
  clusterIP: None
  selector:
    app: kole-controller 

---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  annotations:
  name: kole-controller
  namespace: kole 
spec:
  replicas: 0 
  selector:
    matchLabels:
      app: kole-controller
  serviceName: kole-controller
  template:
    metadata:
      labels:
        app: kole-controller
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: nodepool
                operator: In
                values:
                - kole-controller
      containers:
      - args:
        - --v=2
        command:
        - /kole-controller
        env:
        - name: MQTT5_SERVER 
          value: "${MQTT5_SERVER}"
        - name: HB_TIMEOUT
          value: "300"
        - name: NAME_SPACE
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
        image: ${IMAGES} 
        imagePullPolicy: Always
        name: kole-controller
        resources:
          limits:
            cpu: "8"
            memory: 32Gi
          requests:
            cpu: "1"
            memory: 2Gi
        volumeMounts:
        - mountPath: /etc/localtime
          name: volume-localtime
      dnsPolicy: ClusterFirst
      hostNetwork: true
      restartPolicy: Always
      tolerations:
      - operator: Exists
      volumes:
      - hostPath:
          path: /etc/localtime
        name: volume-localtime

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: lite-kubelet-start-signal
  namespace: kole
data:
  test: start
---
apiVersion: v1
kind: Service
metadata:
  name: lite-kubelet 
  namespace: kole 
  labels:
    app: lite-kubelet
spec:
  ports:
  - port: 80
    name: lite-kubelet
  clusterIP: None
  selector:
    app: lite-kubelet

---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: lite-kubelet
  namespace: kole 
spec:
  replicas: 0
  selector:
    matchLabels:
      app: lite-kubelet
  serviceName: lite-kubelet
  template:
    metadata:
      labels:
        app: lite-kubelet
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: nodepool
                operator: In
                values:
                - lite-kubelet 
      containers:
      - args:
        - --v=2
        command:
        - /lite-kubelet
        env:
        - name: SIGNAL_CM_NAME 
          value: "lite-kubelet-start-signal"
        - name: SIMULATIONS_NUMS
          value: "5000"
        - name: CREATE_CLIENT_INTERVAL
          value: "10"
        - name: HEAT_BEAT_INTERVAL
          value: "120"
        - name: MQTT5_SERVER 
          value: "${MQTT5_SERVER}"
        - name: POD_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        - name: NAME_SPACE
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
        image: ${IMAGES} 
        imagePullPolicy: Always
        name: lite-kubelet
        resources:
          limits:
            cpu: "8"
            memory: 10000Mi
          requests:
            cpu: "2"
            memory: 4000Mi
        volumeMounts:
        - mountPath: /etc/localtime
          name: volume-localtime
      dnsPolicy: ClusterFirst
      hostNetwork: true
      restartPolicy: Always
      schedulerName: default-scheduler
      tolerations:
      - operator: Exists
      volumes:
      - hostPath:
          path: /etc/localtime
          type: ""
        name: volume-localtime
EOF
}

function create_yaml_file() {
    dir=$1
    for file in $dir/*
    do
        file_name=$(basename $file)
        if ! test -f $file ; then
            #echo "$file_name is not a file"
            continue
        fi
        extension="${file##*.}"
        if [ "$extension" != "yaml" ]; then
            #echo "$file_name is not a yaml file"
            continue
        fi
        cat "$file" >> $YAML_FILE
        echo "---" >> $YAML_FILE
    done
}

function manifest() {
    create_basefile
    create_yaml_file $CONFIG_CRD_DIR
    create_yaml_file $CONFIG_RBAC_DIR
}

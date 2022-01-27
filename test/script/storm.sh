# Copyright 2022 The OpenYurt Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#!/bin/bash

BINARY=$0
SUB_CMD=$1

KOLE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

TEE="tee -a /tmp/storm_test.log"


###################### Pod nums #######################
POD_NUMS=`expr 5`

#TOTAL_LITE_KUBELET_NUM_POOLS=(100000 200000 300000 400000)
TOTAL_LITE_KUBELET_NUM_POOLS=(100000)

TOTAL_LITE_KUBELET_NUM=`expr 100`

# The time interval for starting an MQTT CLIENT is in ms
CREATE_CLIENT_INTERVAL=`expr 10`

# Lite-kubelet Heartbeat interval, in seconds
LITE_KUBELET_HEAT_BEAT_INTERVAL=$[60 * 2]

# Heartbeat timeout period of the kole-controller (s)
CONTROLLER_HB_TIMEOUT=300

############################################################
GREP_CREATE_CLIENT="Create client and Subscribe all topic successful node"
GREP_REGISTER_NODE="registering successful node"
LITE_NAME="lite-kubelet"
CONTROLLER_NAME="kole-controller"
NAME_SPACE="kole"

KUBECTL="kubectl -n $NAME_SPACE"

# Need to calculate according to TOTAL_LITE_KUBELET_NUM and POD_NUMS
SIMULATIONS_NUMS=`expr $TOTAL_LITE_KUBELET_NUM / $POD_NUMS`


function restartTest() {
    echo "############# start test #################" | $TEE

    $KUBECTL scale statefulset $LITE_NAME --replicas 0
    
    $KUBECTL delete summaries.lite.openyurt.io --all

    CONTROLLER_PATCH=$(cat <<- EOF
spec:
  template:
    spec:
      containers:
      - name: kole-controller 
        env:
        - name: HB_TIMEOUT 
          value: "$CONTROLLER_HB_TIMEOUT"
EOF
)

    $KUBECTL patch statefulsets.apps $CONTROLLER_NAME --patch "$CONTROLLER_PATCH"
    
    $KUBECTL scale statefulset $CONTROLLER_NAME --replicas 1 
    
    echo "The total number of MQTT clients to be created is $TOTAL_LITE_KUBELET_NUM" | $TEE

    echo "The total number of pods to be created is $POD_NUMS" | $TEE

    echo "The number of connections simulated by each POD is $SIMULATIONS_NUMS" | $TEE
    echo "The connection creation interval for each POD simulation is $CREATE_CLIENT_INTERVAL ms" | $TEE
    echo "The heartbeat timeout of $CONTROLLER_NAME is $CONTROLLER_HB_TIMEOUT s" | $TEE

    LITE_KUBELET_PATCH=$(cat <<- EOF
spec:
  template:
    spec:
      containers:
      - name: lite-kubelet
        env:
        - name: SIMULATIONS_NUMS
          value: "$SIMULATIONS_NUMS"
        - name: CREATE_CLIENT_INTERVAL
          value: "$CREATE_CLIENT_INTERVAL"
        - name: HEAT_BEAT_INTERVAL
          value: "$LITE_KUBELET_HEAT_BEAT_INTERVAL"
EOF
)

    $KUBECTL patch statefulsets.apps ${LITE_NAME} --patch "$LITE_KUBELET_PATCH"

    sleep 10 
    $KUBECTL scale statefulset $LITE_NAME --replicas $POD_NUMS 
    
    
    sleep 5
    checkAllRegisterSuccess
}


function checkAllRegisterSuccess() {
    echo "---- ----" | $TEE
    local START_TIME=`expr $(date +%s)`
    echo "Simulation lite-kubelet has officially started $(date "+%Y-%m-%d %H:%M:%S") The time stamp is $START_TIME" | $TEE

    local LASTER_LITE_POD_NAME="$LITE_NAME-`expr $POD_NUMS - 1`"

    while :
    do
        $KUBECTL get pod $LASTER_LITE_POD_NAME
        if [ $? -eq 0 ]; then
            echo "最后一个创建成功"
            break 
        fi
        sleep 1 
    done

    local POD_FINISH_TIME=`expr $(date +%s)`
    echo "The time required to create the last POD is `expr $POD_FINISH_TIME - $START_TIME` s" | $TEE

    sleep 5 
    
    local GREP_DATA="All lite-kubelet registering successfull"
    MAX_TIME=`expr 0` 
    while :
    do
        local REGISTERD_COUNT=`expr 0`
        for ((i=0; i<=`expr $POD_NUMS - 1`; i ++));
        do
            P_N="$LITE_NAME-$i"
            REGISTER_DATA=$($KUBECTL logs $P_N | grep "$GREP_DATA")
            if [ $? -eq 0 ]; then
                REGISTERD_COUNT=`expr $REGISTERD_COUNT + 1`
                D=$(echo $REGISTER_DATA |awk -F '##' '{print $2}')
                R_TIME=`expr $D`
                if [ $R_TIME -gt $MAX_TIME ]; then
                    MAX_TIME=$R_TIME
                fi
                echo "$P_N registerd time $R_TIME"
                echo "max registerd time $MAX_TIME"
            fi
        done
        if [ "$REGISTERD_COUNT" == "$POD_NUMS" ]; then
            break
        fi
        sleep 1
    done
    
    local LITE_FINISH_TIME=`expr $(date +%s)`
    echo "All Lite-Kubelet registrations take time to complete is $MAX_TIME ms" | $TEE

}

function checkWorkLoad() {
    INFDAEMON_SET_NAME="daemon-set1"
    $KUBECTL delete infdaemonsets.lite.openyurt.io $INFDAEMON_SET_NAME 

    echo "---- Test Workload Distribution ---" | $TEE
    local START_TIME=`expr $(date +%s)`
    echo "Simulated Workload distribution starts at $(date "+%Y-%m-%d %H:%M:%S") The time stamp is $START_TIME" | $TEE
    cat <<EOF | kubectl apply -f -
apiVersion: lite.openyurt.io/v1alpha1
kind: InfDaemonSet
metadata:
  name: "${INFDAEMON_SET_NAME}"
  namespace: "$NAME_SPACE"
spec:
  image: "nginx:v1.19.0"
EOF

    local GREP_DATA="All lite-kubelet receive pod data"
    MAX_TIME=`expr 0` 
    while :
    do
        local RECEIVE_COUNT=`expr 0`
        for ((i=0; i<=`expr $POD_NUMS - 1`; i ++));
        do
            P_N="$LITE_NAME-$i"
            REGISTER_DATA=$($KUBECTL logs $P_N | grep "$GREP_DATA")
            if [ $? -eq 0 ]; then
                RECEIVE_COUNT=`expr $RECEIVE_COUNT + 1`
                D=$(echo $REGISTER_DATA |awk -F '##' '{print $2}')
                R_TIME=`expr $D`
                if [ $R_TIME -gt $MAX_TIME ]; then
                    MAX_TIME=$R_TIME
                fi
                echo "$P_N receive time $R_TIME"
                echo "max receive time $MAX_TIME"
            fi
        done
        if [ "$RECEIVE_COUNT" == "$POD_NUMS" ]; then
            break
        fi
    done
    
    local LITE_FINISH_TIME=`expr $(date +%s)`
    echo "Pod distribution to all Lite-Kubelet takes time is  `expr $MAX_TIME - $START_TIME` s" | $TEE

    $KUBECTL delete infdaemonsets.lite.openyurt.io $INFDAEMON_SET_NAME 
}

function checkSimpleWorkLoad() {

    echo "---- Simple rotation test workload distribution ----" | $TEE

    local START_TIME=`expr $(date +%s)`
    echo "Simulated Workload distribution starts at $(date "+%Y-%m-%d %H:%M:%S") the timestamp is $START_TIME" | $TEE

    local GREP_DATA="All lite-kubelet receive pod data"
    MAX_TIME=`expr 0` 
    while :
    do
        local RECEIVE_COUNT=`expr 0`
        for ((i=0; i<=`expr $POD_NUMS - 1`; i ++));
        do
            P_N="$LITE_NAME-$i"
            REGISTER_DATA=$($KUBECTL logs $P_N | grep "$GREP_DATA")
            if [ $? -eq 0 ]; then
                RECEIVE_COUNT=`expr $RECEIVE_COUNT + 1`
                D=$(echo $REGISTER_DATA |awk -F '##' '{print $2}')
                R_TIME=`expr $D`
                if [ $R_TIME -gt $MAX_TIME ]; then
                    MAX_TIME=$R_TIME
                fi
                echo "$P_N receive time $R_TIME"
                echo "max receive time $MAX_TIME"
            fi
        done
        if [ $RECEIVE_COUNT -gt `expr $POD_NUMS - 3` ]; then
            break
        fi
    done
    
    local LITE_FINISH_TIME=`expr $(date +%s)`
    echo "Pod distribution to all Lite-Kubelet requires a maximum time is $MAX_TIME s" | $TEE
}

function cleanAll() {

    echo "scale ${LITE_NAME} replicas to 0"
    $KUBECTL scale statefulset ${LITE_NAME} --replicas 0
    
    echo "delete all pod of ${LITE_NAME} statefulset"
    $KUBECTL delete pod --all
    
    
    echo "scale $CONTROLLER_NAME replicas to 0"
    $KUBECTL scale statefulset $CONTROLLER_NAME --replicas 0

    echo "delete all pod of $CONTROLLER_NAME statefulset"
    $KUBECTL delete pod --all
    
    echo "delete all summaries in $NAME_SPACE ns"
    $KUBECTL delete summaries.lite.openyurt.io --all
    
    sleep 1
}


function rangelog() {
    GREP_NUM="\-`expr $SIMULATIONS_NUMS - 1`"
    
    pods=$($KUBECTL get pod --sort-by=.metadata.creationTimestamp |grep ${LITE_NAME} |awk -F ' ' '{print $1}')
    
    for p in $pods ; do
    
        echo ""
    
        echo "$p"
        $KUBECTL logs $p #| grep  "Registering successful num is $SIMULATIONS_NUMS"
    done
}


function checkRegister() {
    
    GREP_NUM="\-`expr $SIMULATIONS_NUMS - 1`"
    
    pods=$($KUBECTL get pod --sort-by=.metadata.creationTimestamp |grep ${LITE_NAME} |awk -F ' ' '{print $1}')
    
    for p in $pods ; do
    
        echo ""
    
        echo "$p"
        $KUBECTL logs $p  |grep "${GREP_REGISTER_NODE}"| grep "${LITE_NAME}-" |grep "$GREP_NUM" 
    
        if [ $? -ne 0 ]; then
            $KUBECTL logs $p  |grep "${GREP_REGISTER_NODE}"| grep "${LITE_NAME}-" |grep "$GREP_NUM" 
            echo "$p not all registrations are complete"
            echo ""
        fi
    done

}


function startLoop() {

    echo "" | $TEE
    echo "" | $TEE
    echo "############# Start A New Round OF Full Test  #################" | $TEE

    for element in ${TOTAL_LITE_KUBELET_NUM_POOLS[@]}; do
        TOTAL_LITE_KUBELET_NUM=`expr $element`
        SIMULATIONS_NUMS=`expr $TOTAL_LITE_KUBELET_NUM / $POD_NUMS`
        cleanAll 
        restartTest
        sleep 70 
    done
}


function help() {
        echo """
$BINARY start
    Clear all test environments first, then restart the test
$BINARY clean 
    Clear all test environments
$BINARY checkregister
    Test whether all hosts are registered successfully
$BINARY loop 
    Tests are cycled according to different node levels
$BINARY rangelog 
    Loop through all lite-Kubelet pods logs
$BINARY checkworkload 
    Iterate through all lite-Kubelet logs to determine whether the POD is delivered. This operation requires creating infdaemonset resources
$BINARY checkSimpleWorkLoad 
    Loop through all lite-Kubelet logs to check whether the POD is delivered. This operation does not require creating infdaemonset resources

        """

        exit 0
}

case $SUB_CMD in
    "start")
        cleanAll
        restartTest 
        ;;
    "clean")
        cleanAll 
        ;;
    "checkregister")
        checkRegister
        ;;
    "loop")
        startLoop 
        ;;
    "rangelog")
        rangelog 
        ;;
    "checkworkload")
        checkWorkLoad 
        ;;
    "checkSimpleWorkLoad")
        checkSimpleWorkLoad 
        ;;
    *)
        echo "wrong cmd , help info:"
        help
esac

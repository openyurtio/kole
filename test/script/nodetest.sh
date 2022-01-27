#!/bin/bash

BINARY=$0
SUB_CMD=$1

KOLE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

BUILD_BIN="$KOLE_ROOT/_output/bin/stormtest"

############## 不可修改 #############

DEFAULT_ADD_NUMS=`expr 100000`

LIST_PATCH=500
CREATE_BATCH=10
DELETE_BATCH=500

#################
function deleteNode() {
   echo "-----------------  首先清理环境 分批清理个数 $DELETE_BATCH -------------" |tee -a ./list.log
   $BUILD_BIN delete infedgenode --delete-patch-num=$DELETE_BATCH >>list.log 2>&1
   echo "---------------- 环境清理成功 ----------------" |tee -a ./list.log
   exit 0
}

function addNode() {

    echo "########## 开始测试 新赠 $DEFAULT_ADD_NUMS 每次批量创建 $CREATE_BATCH ##########" |tee -a ./list.log
    
    local START_TIME=`expr $(date +%s)`
    echo "开始创建node ..."  |tee -a ./list.log
    $BUILD_BIN create infedgenode --node-num $DEFAULT_ADD_NUMS --batch-num=$CREATE_BATCH --is-small-size=false >>list.log 2>&1
            
    CREATE_END=`expr $(date +%s)`
    echo "创建$DEFAULT_ADD_NUMS node 耗时: `expr $CREATE_END - $START_TIME` s"  |tee -a ./list.log


    echo "开始测试加载时间(分页数为$LIST_PATCH)" |tee -a ./list.log
    $BUILD_BIN load infedgenode --patch-num=$LIST_PATCH >>list.log 2>&1
    if [ $? -ne 0 ]; then
        echo "测试load(分页数$LIST_PATCH) 失败" |tee -a ./list.log
        exit 1 
    fi
    
    NOW_END=$(date "+%Y-%m-%d %H:%M:%S")
    echo ""
    echo "总体测试成功，测试完成时间: $NOW_END" |tee -a ./list.log
    return 0
}

function testss() {
	return 0
}

function addNodeLoop() {

    echo "########## 开始Loop测试 每次新赠 $DEFAULT_ADD_NUMS 每次批量创建 $CREATE_BATCH ##########" |tee -a ./list.log
    for ((i=0; i<=`expr 10`; i ++));
    do
        echo "第${i}次新赠##########" |tee -a ./list.log
        addNode
        if [ $? -ne 0 ]; then
            echo "第${i}次新赠失败" |tee -a ./list.log
            break
        fi
        sleep 5
    done
    
    echo "---------- Loop 测试完成 ----------" |tee -a ./list.log
}

function help() {
        echo """
$BINARY add 
    增加InfEdgeNode
$BINARY delete 
    清空所有InfEdgeNode
$BINARY addloop 
    循环增加InfEdgeNode

        """

        exit 0
}

case $SUB_CMD in
    "add")
        addNode 
        ;;
    "addloop")
        addNodeLoop 
        ;;
    "delete")
        deleteNode 
        ;;
    *)
        echo "wrong cmd , help info:"
        help
esac


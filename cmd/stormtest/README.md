

## 构建二进制 

 首先执行 ./build.sh ，在本地环境构建 出  stormtest 可执行文件
```
./build.sh
```



## 测试说明：

本测试对象 为 `infedgenodes.lite.openyurt.io` 资源， 而且全部在 summarystorm 命名空间下

可以使用下面命令查看
```
kubectl get infedgenodes.lite.openyurt.io -n summarystorm

```


## stormtest 说明
```

Usage:
  stormtest [command]

Available Commands:
  createnodes 在特定命名空间下创建一定数量的  infedgenodes.lite.openyurt.io  实例
  deletenodes 在特定命名空间下删除所有的infedgenodes.lite.openyurt.io  实例
  help        Help about any command
  loadnodes   计算在特定命名空间下加载所有的infedgenodes.lite.openyurt.io  实例 所需要的时间（ms）

Flags:
  -h, --help                 help for stormtest
      --kube-config string   config file (default is $HOME/.kube/config) (default "/Users/bingyu/.kube/config")

Use "stormtest [command] --help" for more information about a command.
```

### stormtest createnodes 

```
./stormtest createnodes -h

在特定命名空间下创建一定数量的  infedgenodes.lite.openyurt.io  实例

Usage:
  stormtest createnodes [flags]

Flags:
  -h, --help           help for createnodes
      --node-num int   the number of nodes (default 10)

Global Flags:
      --kube-config string   config file (default is $HOME/.kube/config) (default "/Users/bingyu/.kube/config")

```

有个参数
 + --node-num  代表了 要创建多少个 infedgenodes.lite.openyurt.io 实例

### stormtest deletenodes

```
./stormtest deletenodes -h


在特定命名空间下删除所有的infedgenodes.lite.openyurt.io  实例

Usage:
  stormtest deletenodes [flags]

Flags:
  -h, --help   help for deletenodes

Global Flags:
      --kube-config string   config file (default is $HOME/.kube/config) (default "/Users/bingyu/.kube/config")

```

表示会将 summarystorm 命名空间下的所有  infedgenodes.lite.openyurt.io 对象都删掉。 相当于测试环境清空。所以每次做测试时， 都应该先执行deletenodes 




### stormtest loadnodes


```
计算在特定命名空间下加载所有的infedgenodes.lite.openyurt.io  实例 所需要的时间（ms）

Usage:
  stormtest loadnodes [flags]

Flags:
  -h, --help   help for loadnodes

Global Flags:
      --kube-config string   config file (default is $HOME/.kube/config) (default "/Users/bingyu/.kube/config")

```

表示计算 summarystorm 命名空间下的加载玩所有  infedgenodes.lite.openyurt.io 对象所需要的时间(ms),
 

## 测试步骤 

若要测试1k (对应stormtest createnodes 的--node-num 参数 为1000) 基本的 infedgenodes

+ 使用`stormtest deletenodes` 命令删除 summarystorm ns 下的所有infedgenodes 资源对象
+ 使用` stormtest createnodes --node-num 1000` 创建 对应的 infedgenodes 对象
+ 使用 `stormtest loadnodes` 命令来计算 完全list 刚刚创建1000 个node，需要的时间


若要测试 10k node

+ 使用`stormtest deletenodes` 命令删除 summarystorm ns 下的所有infedgenodes 资源对象
+ 使用` stormtest createnodes --node-num 10000` 创建 对应的 infedgenodes 对象(数量是10000)
+ 使用 `stormtest loadnodes` 命令来计算 完全list 10000 个node所需要的时间




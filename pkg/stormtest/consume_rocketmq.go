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
	"context"
	"strings"
	"time"

	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/gogap/errors"
	"k8s.io/klog/v2"

	"github.com/openyurtio/kole/cmd/stormtest/cmd/options"
	"github.com/openyurtio/kole/pkg/data"
)

type ConsumeRocketMQ struct {
	StopCh     chan struct{}
	MQConsumer mq_http_sdk.MQConsumer
	MessageNum int
}

func NewConsumeRocketMQ(config *options.RocketMQFlags) (*ConsumeRocketMQ, error) {
	// 设置HTTP协议客户端接入点，进入消息队列RocketMQ版控制台实例详情页面的接入点区域查看。
	//endpoint := config.Endpoint
	// AccessKey ID阿里云身份验证，在阿里云RAM控制台创建。
	//accessKey := "${ACCESS_KEY}"
	// AccessKey Secret阿里云身份验证，在阿里云RAM控制台创建。
	//secretKey := "${SECRET_KEY}"
	// 消息所属的Topic，在消息队列RocketMQ版控制台创建。
	//不同消息类型的Topic不能混用，例如普通消息的Topic只能用于收发普通消息，不能用于收发其他类型的消息。
	topic := "storm"
	// Topic所属的实例ID，在消息队列RocketMQ版控制台创建。
	// 若实例有命名空间，则实例ID必须传入；若实例无命名空间，则实例ID传入null空值或字符串空值。实例的命名空间可以在消息队列RocketMQ版控制台的实例详情页面查看。
	//instanceId := "${INSTANCE_ID}"
	// 您在控制台创建的Group ID。
	//groupId := "${GROUP_ID}"

	client := mq_http_sdk.NewAliyunMQClient(config.Endpoint, config.AccessKey, config.AccessSecret, "")

	mqConsumer := client.GetConsumer(config.Instance, topic, config.Group, "")

	t := &ConsumeRocketMQ{
		MQConsumer: mqConsumer,
		StopCh:     make(chan struct{}, 1),
		MessageNum: config.MessageNum,
	}

	klog.V(4).Infof("create rocket mq client successful")

	return t, nil
}

func (n *ConsumeRocketMQ) ComsumeRocketMq(message string) {
	hb, err := data.UnmarshalPayloadToHeartBeat([]byte(message))
	if err != nil {
		klog.Errorf("UnmarshalPayloadToHeartBeat error %v", err)
		return
	}
	klog.Infof("Receive Message %s", hb.Name)
}

func (n *ConsumeRocketMQ) Run(ctx context.Context) error {
	index := 0
	for {
		endChan := make(chan int)
		respChan := make(chan mq_http_sdk.ConsumeMessageResponse)
		errChan := make(chan error)
		if index >= n.MessageNum {
			klog.Warningf("Received %d message , return loop", index)
			return nil
		}
		go func() {
			select {
			case <-n.StopCh:
				klog.Warningf("Receive stop single, stop")
				return
			case <-ctx.Done():
				klog.Warningf("Receive ctx cancel, stop")
				return
			case resp := <-respChan:
				{
					// 处理业务逻辑。
					var handles []string
					klog.V(5).Infof("Consume %d messages---->\n", len(resp.Messages))
					for _, v := range resp.Messages {
						handles = append(handles, v.ReceiptHandle)
						klog.V(5).Infof("\tMessageID: %s, PublishTime: %d, MessageTag: %s\n"+
							"\tConsumedTimes: %d, FirstConsumeTime: %d, NextConsumeTime: %d\n"+
							"\tBody: %s\n"+
							"\tProps: %s\n",
							v.MessageId, v.PublishTime, v.MessageTag, v.ConsumedTimes,
							v.FirstConsumeTime, v.NextConsumeTime, v.MessageBody, v.Properties)
						n.ComsumeRocketMq(v.MessageBody)
						index++
						if index >= n.MessageNum {
							break
						}
					}

					// NextConsumeTime前若不确认消息消费成功，则消息会被重复消费。
					// 消息句柄有时间戳，同一条消息每次消费拿到的都不一样。
					ackerr := n.MQConsumer.AckMessage(handles)
					if ackerr != nil {
						// 某些消息的句柄可能超时，会导致消息消费状态确认不成功。
						klog.Errorf("Ack message error %v", ackerr)
						if errAckItems, ok := ackerr.(errors.ErrCode).Context()["Detail"].([]mq_http_sdk.ErrAckItem); ok {
							for _, errAckItem := range errAckItems {
								klog.Errorf("\tErrorHandle:%s, ErrorCode:%s, ErrorMsg:%s\n",
									errAckItem.ErrorHandle, errAckItem.ErrorCode, errAckItem.ErrorMsg)
							}
						} else {
							klog.Warningf("ack err =", ackerr)
						}
						time.Sleep(time.Duration(3) * time.Second)
					} else {
						klog.V(4).Infof("Ack ---->\n\t%s\n", handles)
					}

					endChan <- 1
				}
			case err := <-errChan:
				{
					// Topic中没有消息可消费。
					if strings.Contains(err.(errors.ErrCode).Error(), "MessageNotExist") {
						klog.Warningf("\nNo new message, continue!")
					} else {
						klog.Errorf("erroChan %v", err)
						time.Sleep(time.Duration(3) * time.Second)
					}
					endChan <- 1
				}
			case <-time.After(35 * time.Second):
				{
					klog.Warningf("Timeout of consumer message ??")
					endChan <- 1
				}
			}
		}()

		// 长轮询消费消息，网络超时时间默认为35s。
		// 长轮询表示如果Topic没有消息，则客户端请求会在服务端挂起3s，3s内如果有消息可以消费则立即返回响应。
		n.MQConsumer.ConsumeMessage(respChan, errChan,
			10, // 一次最多消费3条（最多可设置为16条）。
			3,  // 长轮询时间3s（最多可设置为30s）。
		)
		<-endChan
	}

}

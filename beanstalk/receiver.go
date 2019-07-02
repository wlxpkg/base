/*
 * @Author: qiuling
 * @Date: 2019-07-02 15:59:10
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-02 17:25:46
 */
package beanstalk

import (
	. "artifact/pkg"
)

type Receiver interface {
	OnReceive(string) bool // 处理收到的消息
}

type Receive struct {
}

func NewReceiver() *Receive {
	r := &Receive{}

	return r
}

// OnReceive 消息处理, 此处为demo
func (r *Receive) OnReceive(jsonStr string) bool {
	R(jsonStr, "收到 beanstalk 消息")
	json, _ := JsonDecode(jsonStr)

	topic := json["topic"].(string)
	R(topic, "topic")

	message := json["message"].(map[string]interface{})
	R(message, "message")

	msgName := message["name"].(string)
	R(msgName, "msgName")

	return true
}

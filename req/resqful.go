/*
 * @Author: qiuling
 * @Date: 2019-06-25 20:44:57
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-26 18:16:33
 */
package req

import (
	. "artifact/pkg"
	. "artifact/pkg/config"
	"artifact/pkg/log"
	"encoding/json"
	"errors"
)

type Restful struct {
	client  *HttpClient
	exp     bool
	service map[string]string
}

func NewRestful(name string) *Restful {
	r := new(Restful)
	r.exp = true
	r.client = NewClient("")
	r.setService()
	r.GetService(name)
	return r
}

// setService 初始化时候设置服务map
func (r *Restful) setService() {
	config := Config.Server

	serverUrl := make(map[string]string)

	serverUrl["tools"] = "http://" + config.Tools
	serverUrl["user"] = "http://" + config.User
	serverUrl["course"] = "http://" + config.Course
	serverUrl["discovery"] = "http://" + config.Discovery
	serverUrl["common"] = "http://" + config.Common
	serverUrl["grant"] = "http://" + config.Grant
	serverUrl["shop"] = "http://" + config.Shop
	serverUrl["message"] = "http://" + config.Message

	r.service = serverUrl
}

// GetService 获取一个服务地址设置给 http 客户端
func (r *Restful) GetService(name string) *Restful {
	baseUrl, exists := r.service[name]

	if !exists {
		log.Err("Restful.GetService 服务不存在, name:" + name)
	}
	r.client.SetBaseUrl(baseUrl)
	return r
}

// SetJwt 设置 jwt
func (r *Restful) SetJwt(jwt string) *Restful {
	r.client.AddHeader("Authorization", "Bearer "+jwt)
	return r
}

// 设置 exp, 为 true 则不返回原始数据直接异常
// 默认 true
func (r *Restful) SetExp(exp bool) *Restful {
	r.exp = exp
	return r
}

func (r *Restful) SetData(data map[string]string) *Restful {
	r.client.SetData(data)
	return r
}

// Req 发送请求
func (r *Restful) Req(method string, route string) (data interface{}, err error) {
	resp, err := r.client.Request(method, route)

	if err != nil {
		dataStr, _ := json.Marshal(r.client.data)
		log.Warn("微服务请求失败! service: " + r.client.baseUrl + " data: " + Byte2String(dataStr) + "method: " + method + "route: " + route)
		return
	}

	data, err = r.serviceData(resp)
	return
}

// serviceData 解析数据
func (r *Restful) serviceData(resp string) (resData interface{}, err error) {
	// R(resp, "resp")
	var data RespData
	err = json.Unmarshal([]byte(resp), &data)

	if err != nil {
		dataStr, _ := json.Marshal(r.client.data)
		log.Warn("微服务数据解析失败! service: " + r.client.baseUrl + " req: " + Byte2String(dataStr) + "resp: " + resp + "err:" + err.Error())
		return
	}

	if r.exp {
		if data.Code != 1 {
			err = errors.New(data.Message)
			return
		}
		resData = data.Data
	} else {
		resData = data
	}

	return
}

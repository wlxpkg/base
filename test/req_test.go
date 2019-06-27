/*
 * @Author: qiuling
 * @Date: 2019-06-25 17:34:31
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-26 17:19:50
 */
package test

import (
	. "artifact/pkg"
	"artifact/pkg/req"
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = make(map[string]string)
var header = make(map[string]string)

func TestHttp(t *testing.T) {

	client := req.NewClient("http://user.artifact.com")

	// GET
	data["aaa"] = "bbb"
	data["ccc"] = "ddd"
	response, err := client.SetData(data).Request("GET", "/test/osstoken")

	R(response, "response")
	R(err, "err")

	res := "{\"code\":10102,\"message\":\"无效的token\"}"

	assert.Equal(t, res, response, "TestLog")

	// POST JSON
	data = make(map[string]string)

	header["Authorization"] = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiIsImp0aSI6In07d1EyfU1STXkifQ.eyJqdGkiOiJ9O3dRMn1NUk15IiwiaWF0IjoxNTYwMzMxOTk5LCJleHAiOjE1OTE4Njc5OTksImFydGlmYWN0IjoiUHowaVBORm9VNUlhSExtcFpRNHVTNU9PdHVwS2dxY1giLCJ0b2tlbiI6IjIyMzMzMzcwZDMzN2YzM2RmMTI5ZTc0MmY1ZWU3MDBkIn0.E8N9F9-hVTmM8NlNkbmREFUe0x_G_-VhAaNcSE4xMw4"

	data["name"] = "测试角色"
	data["slug"] = "customer"
	data["permission_id"] = "4,10,2,1,16,19,23,27,31,35,39,40,36,32,28,24,20,17,18,21,22,26,25,29,30,33,34,37,38,41,42"
	data["type"] = "99"
	data["is_default"] = "0"

	responsePost, err := client.SetData(data).AddHeaders(header).Request("PUT", "/admin/role/38")

	R(responsePost, "response")
	R(err, "err")

	resPost := "{\"code\":1,\"message\":\"\",\"data\":38}"

	assert.Equal(t, resPost, responsePost, "TestLog")
}

func TestRestful(t *testing.T) {
	userService := req.NewRestful("user")
	// R(userService, "userService")

	data["phone"] = "13456724507"
	data["password"] = "11111"

	res, err := userService.SetData(data).Req("POST", "/login/phone")

	R(res, "res")
	R(err, "err")

	assert.Equal(t, 38, res, "TestLog")
}

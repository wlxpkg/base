/*
 * @Author: qiuling
 * @Date: 2019-05-10 14:43:57
 * @Last Modified by: qiuling
 * @Last Modified time: 2020-07-07 21:58:39
 */

// http 请求连接池
package req

import (
	"bytes"

	"github.com/wlxpkg/base/log"

	jsoniter "github.com/json-iterator/go"

	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// var client *http.Client
const (
	MaxIdleConns        = 100
	MaxIdleConnsPerHost = 100
	IdleConnTimeout     = 90
	Timeout             = 30
)

type HttpClient struct {
	BaseURL string
	Header  map[string]string
	Data    map[string]string
	Client  *http.Client
}

// init HTTPClient
func NewClient(baseUrl string) *HttpClient {
	c := new(HttpClient)

	c.BaseURL = baseUrl
	c.Header = make(map[string]string)
	c.Data = make(map[string]string)
	c.Client = createHTTPClient()
	return c
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   Timeout * time.Second,
				KeepAlive: Timeout * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     IdleConnTimeout * time.Second,
		},
	}
	return client
}

// SetBaseURL 设置基础地址
func (c *HttpClient) SetBaseURL(baseURL string) *HttpClient {
	c.BaseURL = baseURL
	return c
}

// AddHeader 添加一个 header
func (c *HttpClient) AddHeader(name string, value string) *HttpClient {
	c.Header[name] = value
	return c
}

// AddHeaders 添加一组 header 数据
func (c *HttpClient) AddHeaders(headers map[string]string) *HttpClient {
	for key, header := range headers {
		c.Header[key] = header
	}
	return c
}

// SetData 设置查询数据
func (c *HttpClient) SetData(data map[string]string) *HttpClient {
	c.Data = data
	return c
}

// reset 重置 header 和 data
func (c *HttpClient) reset() {
	c.Header = make(map[string]string)
	c.Data = make(map[string]string)
}

// Request 发送 http 请求
func (c *HttpClient) Request(method string, route string) (resp string, err error) {

	if c.BaseURL == "" {
		err = errors.New("ERR_HTTP_BASEURL")
		return
	}

	// 这个地方有点恶心😖...算是个坑吧
	postType := method
	if method == "FORM" {
		method = "POST"
	}

	requrl, body, err := c.buildData(postType, route)
	if err != nil {
		log.Err(err)
		err = errors.New("ERR_PARAM")
		return
	}

	reqest, err := http.NewRequest(method, requrl, body)
	if err != nil {
		log.Err(err)
		err = errors.New("ERR_HTTP_TIMEOUT")
		return
	}
	reqest = c.buildHeader(reqest, postType)

	response, err := c.Client.Do(reqest) //发送请求
	c.reset()
	if err != nil {
		log.Err(err)
		err = errors.New("ERR_HTTP_TIMEOUT")
		return
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Err(err)
		err = errors.New("ERR_HTTP_TIMEOUT")
		return
	}

	resp = string(data)
	return
}

// buildData 构建请求参数
// 除了 GET 拼到url, FORM post 表单, 默认 body json string
// method: GET FORM( POST FORM) POST PUT DELETE PATCH HEAD OPTIONS
func (c *HttpClient) buildData(method string, route string) (requrl string, body *strings.Reader, err error) {
	body = strings.NewReader("")
	requrl = c.BaseURL + route
	if len(c.Data) != 0 {
		switch method {
		case "GET":
			query := url.Values{}
			for k, v := range c.Data {
				query.Add(k, v)
			}
			requrl += "?" + query.Encode()
		case "FORM":
			v := url.Values{}
			for key, data := range c.Data {
				v.Set(key, data)
			}

			body = strings.NewReader(v.Encode())
		default:
			jsonBype, jsonErr := json.Marshal(c.Data)
			if jsonErr != nil {
				err = jsonErr
				return
			} else {
				jsonBuf := bytes.NewBuffer(jsonBype)
				body = strings.NewReader(jsonBuf.String())
			}
		}
	}
	return
}

func (c *HttpClient) buildHeader(reqest *http.Request, postType string) *http.Request {
	if len(c.Header) != 0 {
		for k, h := range c.Header {
			reqest.Header.Set(k, h)
		}
	}

	if postType == "FORM" {
		reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	} else {
		reqest.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	return reqest
}

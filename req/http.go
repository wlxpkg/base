/*
 * @Author: qiuling
 * @Date: 2019-05-10 14:43:57
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-28 18:17:52
 */

// http 请求连接池
package req

import (
	"zwyd/pkg/log"
	"bytes"

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
	baseUrl string
	header  map[string]string
	data    map[string]string
	client  *http.Client
}

// init HTTPClient
func NewClient(baseUrl string) *HttpClient {
	c := new(HttpClient)

	c.baseUrl = baseUrl
	c.header = make(map[string]string)
	c.data = make(map[string]string)
	c.client = createHTTPClient()
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

// SetBaseUrl 设置基础地址
func (c *HttpClient) SetBaseUrl(baseUrl string) *HttpClient {
	c.baseUrl = baseUrl
	return c
}

// AddHeader 添加一个 header
func (c *HttpClient) AddHeader(name string, value string) *HttpClient {
	c.header[name] = value
	return c
}

// AddHeaders 添加一组 header 数据
func (c *HttpClient) AddHeaders(headers map[string]string) *HttpClient {
	for key, header := range headers {
		c.header[key] = header
	}
	return c
}

// SetData 设置查询数据
func (c *HttpClient) SetData(data map[string]string) *HttpClient {
	c.data = data
	return c
}

// reset 重置 header 和 data
func (c *HttpClient) reset() {
	c.header = make(map[string]string)
	c.data = make(map[string]string)
}

// Request 发送 http 请求
func (c *HttpClient) Request(method string, route string) (resp string, err error) {

	if c.baseUrl == "" {
		err = errors.New("ERR_HTTP_BASEURL")
		return
	}

	requrl, body, err := c.buildData(method, route)
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
	reqest = c.buildHeader(reqest, method)

	response, err := c.client.Do(reqest) //发送请求
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
	v := url.Values{}
	body = strings.NewReader("")
	requrl = c.baseUrl + route
	if len(c.data) != 0 {
		switch method {
		case "GET":
			query := url.Values{}
			for k, v := range c.data {
				query.Add(k, v)
			}
			requrl += "?" + query.Encode()
		case "FORM":
			for key, data := range c.data {
				v.Set(key, data)
			}

			body = strings.NewReader(v.Encode())
		default:
			jsonBype, jsonErr := json.Marshal(c.data)
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

func (c *HttpClient) buildHeader(reqest *http.Request, method string) *http.Request {
	if len(c.header) != 0 {
		for k, h := range c.header {
			reqest.Header.Set(k, h)
		}
	}

	if method == "FORM" {
		reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	} else {
		reqest.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	return reqest
}

/*
// Get http get 请求
func (c *HttpClient) Get(requrl string, data map[string]string) (resp string, err error) {
	response, err := c.client.Get(requrl)
	c.reset()
	if err != nil {
		log.Err(err)
		err = errors.New("ERR_HTTP_TIMEOUT")
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Err(err)
		err = errors.New("ERR_HTTP_TIMEOUT")
		return
	}

	resp = string(body)
	return
}

// Post http post 请求
func (c *HttpClient) Post(requrl string, data map[string]string) (resp string, err error) {

	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}

	response, err := c.client.PostForm(requrl, form)
	c.reset()
	if err != nil {
		log.Err(err)
		err = errors.New("ERR_HTTP_TIMEOUT")
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Err(err)
		err = errors.New("ERR_HTTP_TIMEOUT")
		return
	}

	resp = string(body)
	return
}
*/

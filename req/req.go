/*
 * @Author: qiuling
 * @Date: 2019-05-10 14:43:57
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-05-10 16:10:58
 */

// http 请求连接池
package req

import (
	"artifact/pkg/log"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	client *http.Client
)

// init HTTPClient
func init() {
	client = createHTTPClient()
}

const (
	MaxIdleConns        = 100
	MaxIdleConnsPerHost = 100
	IdleConnTimeout     = 90
	Timeout             = 30
)

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

// Get http get 请求
func Get(requrl string) (resp string, err error) {
	response, err := client.Get(requrl)
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
func Post(requrl string, data map[string]string) (resp string, err error) {

	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}

	response, err := client.PostForm(requrl, form)
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

/*
 * @Author: qiuling
 * @Date: 2019-04-29 19:32:36
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-05-05 10:37:19
 */
package pkg

import (
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"runtime/debug"
	"time"
)

func R(data interface{}, name string) {
	fmt.Printf("%v\n", name)
	fmt.Printf("%+v\n", data)
}

func D(data interface{}) {
	fmt.Printf("%s\n", debug.Stack())
	fmt.Printf("%+v\n", data)
}

// JSONTime format json time field by myself
type JSONTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with Y-m-d H:i:s
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan value of time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// 发送get请求
func RequestGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

// Paginator 分页统一结构
func Paginator(page int, count int, list interface{}) map[string]interface{} {
	data := make(map[string]interface{}, 3)
	data["total_count"] = count
	data["current_page"] = page
	data["list"] = list

	return data
}

func RandomNum(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

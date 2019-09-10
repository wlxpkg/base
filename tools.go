/*
 * @Author: qiuling
 * @Date: 2019-04-29 19:32:36
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-08-22 15:48:23
 */
package pkg

import (
	"artifact/pkg/log"
	"crypto/md5"
	"fmt"
	"math/big"
	"math/rand"
	"runtime/debug"
	"strconv"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"

	"github.com/tidwall/gjson"
	"github.com/tv42/base58"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type RespData struct {
	Code    int64
	Message string
	Data    interface{}
}

const TimeFormat = "2006-01-02 15:04:05"

func R(data interface{}, name string) {
	fmt.Printf("\n"+name+": \n%#v\n", data)
}

func D(data interface{}) {
	fmt.Printf("\n%s :\n", debug.Stack())
	fmt.Printf("\n%#v\n", data)
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

// String2Time 字符串格式[转换为]time,Time(时间对象) eg:(2019-09-09T09:09:09+08:00)
func String2Time(paramTime string) (returnTime time.Time) {

	loc, _ := time.LoadLocation("Asia/Shanghai")
	returnTime, _ = time.ParseInLocation(TimeFormat, paramTime, loc)
	return
}

// String2Unix 字符串格式[转换为]时间戳 eg:(1557398617)
func String2Unix(paramTime string) int64 {
	timeStruct := String2Time(paramTime)
	second := timeStruct.Unix()
	return second
}

// Time2Unix 时间对象[转换为]Unix时间戳 eg:(1557398617)
func Time2Unix(paramTime time.Time) int64 {
	second := paramTime.Unix()
	return second
}

// Time2String 时间对象[转换为]字符串 eg:(2019-09-09 09:09:09)
func Time2String(t time.Time) string {
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(TimeFormat)
	return str
}

// Unix2String 时间戳[转换为]字符串 eg:(2019-09-09 09:09:09)
func Unix2String(stamp int64) string {
	str := time.Unix(stamp, 0).Format(TimeFormat)
	return str
}

// Unix2Time unix时间戳[转换为]时间对象 eg:(2019-09-09T09:09:09+08:00)
func Unix2Time(stamp int64) time.Time {
	stampStr := Unix2String(stamp)
	timer := String2Time(stampStr)
	return timer
}

// 生成随机字符串
func RandStr(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// String2Int 字符串 转 int
func String2Int(str string) (int, error) {
	return strconv.Atoi(str)
}

// String2Int64 字符串 转 int64
func String2Int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// String2Uint string 转 uint
func String2Uint(str string) uint {
	u64, _ := strconv.ParseUint(str, 10, 0)
	return uint(u64)
}

// int 转 string
func Int2String(intval int) string {
	return strconv.Itoa(intval)
}

// int64 转 string
func Int642String(intval int64) string {
	return strconv.FormatInt(intval, 10)
}

func Uint642String(intval uint64) string {
	return strconv.FormatUint(intval, 10)
}

// Byte2String 字节数组转 string
func Byte2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// float64 转 string
func Float642String(v float64) string {
	return strconv.FormatFloat(v, 'g', 20, 64)
}

func Bool2Int(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Data2Map 将数据转为 map
func Data2Map(data interface{}) map[string]interface{} {
	mdata := make(map[string]interface{})
	j, _ := json.Marshal(data)

	err := json.Unmarshal(j, &mdata)
	if err != nil {
		log.Warn(err)
		return mdata
	}
	return mdata
}

// Base58Encode int64 -> 58进制转换
func Base58Encode(id int64) string {
	bigint := new(big.Int).SetInt64(id)

	str := base58.EncodeBig(nil, bigint)
	return Byte2String(str)
}

// Base58Decode 58进制字符串转 int64
func Base58Decode(str string) int64 {
	dint, _ := base58.DecodeToBig([]byte(str))
	return dint.Int64()
}

// Md5
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

// JsonDecode
func JsonEncode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// JsonDecode 你懂的
func JsonDecode(data string) (result map[string]interface{}, err error) {
	result, ok := gjson.Parse(data).Value().(map[string]interface{})
	if !ok {
		log.Warn("数据解析失败! data: " + data)
		err = Excp("ERR_DATA_DECODE")
		return
	}
	return
}

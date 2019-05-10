/*
 * @Author: qiuling
 * @Date: 2019-05-10 14:23:40
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-05-10 19:20:14
 */

package pkg

import (
	"artifact/pkg/log"
	"artifact/pkg/req"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gookit/config"
)

type Model struct {
}

// JSONTime format json time field by myself

// JSONTime format json time field by myself
type JSONTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with Y-m-d H:i:s
func (t JSONTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
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

func CreateID() (idInt int64, err error) {
	domain := config.GetEnv("IDGENERATOR_URL", "http://192.168.3.3")
	uri := domain + "/getid"

	idStr, err := req.Get(uri)
	if err != nil {
		log.Warn(err)
		err = errors.New("ERR_IDGEN_FAIL")
		return
	}

	idInt, _ = strconv.ParseInt(idStr, 10, 64)
	return
}

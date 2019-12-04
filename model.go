/*
 * @Author: qiuling
 * @Date: 2019-05-10 14:23:40
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-17 17:20:16
 */

package pkg

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"git.wlx/zwyd/pkg/cache"
	"git.wlx/zwyd/pkg/log"
	"strings"
	"time"
)

type Model struct {
	// cache *cache.Cache
}

var Cache = cache.NewCache()

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

// CreateID id生成器生成id
func CreateID() (idStr string, err error) {
	idStr, err = cache.Getid()
	if err != nil {
		log.Warn(err)
		err = errors.New("ERR_IDGEN_FAIL")
		return
	}

	return
}

// BatchID id生成器批量生成id
func BatchID(num int) (ids []string, err error) {
	idsStr, err := cache.Batchid(num)
	if err != nil {
		log.Warn(err)
		err = errors.New("ERR_IDGEN_FAIL")
		return
	}

	ids = strings.Split(idsStr, ",")
	return
}

/*
 * @Author: qiuling
 * @Date: 2019-05-09 15:36:17
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-12-05 11:10:17
 */
package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlxpkg/base/log"
)

func TestLog(t *testing.T) {
	err := errors.New("this is a test err")

	log.Info(err)
	log.Err(err)

	assert.Equal(t, 1, 1, "TestLog")
}

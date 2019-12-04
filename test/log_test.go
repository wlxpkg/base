/*
 * @Author: qiuling
 * @Date: 2019-05-09 15:36:17
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-05-09 15:47:10
 */
package test

import (
	"git.wlx/zwyd/pkg/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLog(t *testing.T) {
	err := errors.New("this is a test err")

	log.Info(err)
	log.Err(err)

	assert.Equal(t, 1, 1, "TestLog")
}

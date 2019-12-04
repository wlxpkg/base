/*
 * @Author: qiuling
 * @Date: 2019-05-10 16:46:20
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-05-10 17:09:00
 */
package test

import (
	. "git.wlx/zwyd/pkg"
	"testing"
)

func TestCreateId(t *testing.T) {
	id, _ := CreateID()
	R(id, "CreateID")
}

func BenchmarkCreateId(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, _ := CreateID()
		R(id, "CreateID")
	}
}

package test

import (
	. "artifact/pkg"
	"testing"
)

// BenchmarkRandomNum 随机一定范围内的数字
func BenchmarkRandomNum(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RandomNum(1000, 9999)
	}
}

package test

import (
	"fmt"
	"testing"

	. "github.com/wlxpkg/base"
)

// BenchmarkRandomNum 随机一定范围内的数字
func BenchmarkRandomNum(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RandomNum(1000, 9999)
	}
}

func TestTodayDate(t *testing.T) {
	tm := TodayDate()
	fmt.Println(tm.Format("2006-01-02"))
	fmt.Println(tm.AddDate(0, 0, 1).Format("2006-01-02"))
	fmt.Println(tm.AddDate(0, 0, -1).Format("2006-01-02"))
}

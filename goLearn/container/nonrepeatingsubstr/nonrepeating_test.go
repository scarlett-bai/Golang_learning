package main

import "testing"

func TestSubStr(t *testing.T) {
	tests := []struct {
		s   string
		ans int
	}{
		// Normal cases
		{"abcabcbb", 3},
		{"pwwkew", 3},

		// Edge cases 特殊例子
		{"", 0},
		{"b", 1},
		{"bbbbb", 1},
		{"abcabcabcabcd", 4},

		// Chinese support
		{"这里是慕课网", 6},
		{"一二三二一", 3},
		{"黑化肥挥发发灰会花飞灰化肥会发发黑会飞花", 7},
	}
	for _, tt := range tests {
		actual := lengthOfNonRepeatingSubStr(tt.s)
		if actual != tt.ans {
			t.Errorf("got %d for input %s; expected %d", actual, tt.s, tt.ans)
		}
	}
}

func BenchmarkSubStr(b *testing.B) {
	// 性能测试
	s := "黑化肥挥发发灰会花飞灰化肥会发发黑会飞花"
	ans := 7

	for i := 0; i < b.N; i++ {
		// 既然是性能测试，一定要测很多遍，就设置一个循环，遍历次数我们不用管
		actual := lengthOfNonRepeatingSubStr(s)
		if actual != ans {
			b.Errorf("got %d for input %s; expected %d", actual, s, ans)
		}
	}

}

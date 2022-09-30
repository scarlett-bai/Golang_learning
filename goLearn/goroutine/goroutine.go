package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	var a [10]int
	for i := 0; i < 10; i++ {
		go func() {
			for {
				a[i]++
				runtime.Gosched() // 让协程 手动交出控制权
			}
		}()
	}
	time.Sleep(time.Millisecond * 10)
	fmt.Printf("%v\n", a)
}

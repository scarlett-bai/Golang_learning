package main

import (
	"context"
	"fmt"
	"time"
)

type paramKey struct{}

func main() {
	c := context.WithValue(context.Background(),
		paramKey{}, "abc")
	c, cancel := context.WithTimeout(c, 10*time.Second) // 生成了新的context  然后又赋值给了c
	defer cancel()
	go mainTask(c) // 这是一个总任务

	var cmd string
	for {
		fmt.Scan(&cmd) // 从键盘输入一个值
		if cmd == "c" {
			cancel()
		}
	}
}

func mainTask(c context.Context) {
	fmt.Printf("main task started with param %q\n", c.Value(paramKey{}))
	go func() {
		c1, cancel := context.WithTimeout(c, 10*time.Second) // 将整个context 分出2s的时间给task1
		defer cancel()
		smallTask(c1, "task1", 9*time.Second) // smallTask 不是 mainTask里面的子任务，而是mainTask里面的步骤
	}()

	// go smallTask(context.Background(), "task1", 10*time.Second) // 开一个后台任务，不管你做几秒
	smallTask(c, "task2", 8*time.Second)
}

func smallTask(c context.Context, name string, d time.Duration) {
	fmt.Printf("%s started with param %q\n", name, c.Value(paramKey{}))
	select {
	case <-time.After(d):
		fmt.Printf("%s done\n", name)
	case <-c.Done(): // 是一个channel 收到值
		fmt.Printf("%s cancelled\n", name)
	}
}

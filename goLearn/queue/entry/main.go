package main

import (
	"fmt"
	"goLearn/queue"

	"go.uber.org/zap"
)

func main() {
	q := queue.Queue{1}
	q.Push(2)
	q.Push(3)
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.IsEmpty())
	fmt.Println(q.Pop())
	fmt.Println(q.IsEmpty())
	log, _ := zap.NewProduction()
	log.Warn("This is a Warning logger")
}

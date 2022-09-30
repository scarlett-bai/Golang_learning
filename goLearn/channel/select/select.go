package main

import (
	"fmt"
	"math/rand"
	"time"
)

func worker(id int, c chan int) {
	for n := range c {
		time.Sleep(time.Second)
		fmt.Printf("Worker %d received %d \n", id, n)
		// w.done()

	}
}

// type worker struct {
// 	in   chan int
// 	done func()
// }

func createWorker(id int) chan int {
	c := make(chan int)
	go worker(id, c)
	return c
}

func generator() chan int {
	out := make(chan int)
	go func() {
		i := 0
		for {
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
			out <- i
			i++
			// fmt.Println(i)
		}
	}()
	return out
}

func main() {
	// var c1, c2 chan int // c1 c2 nill
	var c1, c2 = generator(), generator()
	worker := createWorker(0)
	// n := 0
	var values []int
	tm := time.After(time.Second * 10)
	tick := time.Tick(time.Second)
	for {
		var activeWorker chan int
		var activeValue int
		if len(values) > 0 {
			activeWorker = worker
			activeValue = values[0]
			// fmt.Println("activieValue =", activeValue)
		}

		select {
		// 通过select 进行调度
		case n := <-c1:
			values = append(values, n)
		case n := <-c2:
			values = append(values, n)
		case activeWorker <- activeValue:
			values = values[1:]

		// 增加计时器
		case <-time.After(800 * time.Millisecond):
			fmt.Println("Timeout")
		case <-tick:
			fmt.Println("queue length = ", len(values))
		case <-tm:
			fmt.Println("bye")
			return
		}
	}
}

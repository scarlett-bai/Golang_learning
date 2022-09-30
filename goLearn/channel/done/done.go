package main

import (
	"fmt"
	"sync"
)

func doWorker(id int, w worker) {
	for n := range w.in {
		fmt.Printf("Worker %d received %c\n", id, n)
		w.done()

	}
}

type worker struct {
	in   chan int
	done func()
}

func createWorker(id int, wg *sync.WaitGroup) worker {
	w := worker{
		in: make(chan int),
		done: func() {
			wg.Done()
		},
	}
	go doWorker(id, w)
	return w
}

func chanDemo() {
	var wg sync.WaitGroup

	// var c chan int   // c==nil  没有办法用的
	var workers [10]worker
	for i := 0; i < 10; i++ {
		workers[i] = createWorker(i, &wg)
	}
	wg.Add(20)
	for i, worker := range workers {
		worker.in <- 'a' + i
		// <-workers[i].done
	}

	for i, worker := range workers {
		worker.in <- 'A' + i
		// <-workers[i].done
	}

	wg.Wait()
	// wait for all of them
	// for _, worker := range workers {
	// 	<-worker.done
	// 	<-worker.done
	// }
}

func main() {
	chanDemo()
}

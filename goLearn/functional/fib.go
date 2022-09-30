package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func fibonacci() intGen {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func printFileContents(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

// 是一个函数类型  函数类型就可以实现接口
type intGen func() int

func (g intGen) Read(p []byte) (n int, err error) {
	next := g()
	if next > 10000 {
		return 0, io.EOF
	}
	s := fmt.Sprintf("%d\n", next)
	// TODO: incorrect if p is too small
	return strings.NewReader(s).Read(p)
}

func main() {
	// f := fib.Fibonacci()
	// printFileContents()

	// fmt.Println(f())
	// fmt.Println(f())
	// fmt.Println(f())
	// fmt.Println(f())
	// fmt.Println(f())
	// fmt.Println(f())
	// fmt.Println(f())
	// fmt.Println(f())
}

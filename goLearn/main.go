package main

// import (
// 	"fmt"
// 	// "log"
// 	// "time"
// )

// func main() {
// 	// fmt.Println("Hello world!")
// 	// fmt.Println("hello world!!!")
// 	// fmt.Println(time.Now().Unix())
// 	// log.Println("log")
// 	arr := [8]int{0, 1, 2, 3, 4, 5, 6, 7}
// 	s1 := arr[2:6]
// 	fmt.Printf("s1=%v len(s1)=%d cap(s1)=%d\n", s1, len(s1), cap(s1))
// 	s2 := s1[3:5]
// 	fmt.Printf("s1=%v len(s1)=%d cap(s1)=%d\n", s2, len(s2), cap(s2))
// 	s3 := append(s2, 10)
// 	s4 := append(s3, 11)
// 	s5 := append(s4, 12)
// 	fmt.Println("s3, s4, s5 =", s3, s4, s5, s5)
// 	// s4 s5 no longer view arr 但是视图是谁 由系统分配，分配的一个更大的空间存放切片数据
// 	fmt.Println("arr =", arr)
// }

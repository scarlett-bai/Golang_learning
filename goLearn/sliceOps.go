package main

import "fmt"

func printSlice(s []int) {
	fmt.Printf(" %v len=%d, cap=%d\n", s, len(s), cap(s))
}

// func main() {
// 	// var s []int // zero value for slice is nil
// 	// for i := 0; i < 100; i++ {
// 	// 	printSlice(s)
// 	// 	s = append(s, 2*i+1)
// 	// }
// 	// fmt.Println(s)
// 	s1 := []int{2, 4, 6, 8}
// 	s2 := make([]int, 16)
// 	// s3 := make([]int, 10, 32)

// 	fmt.Println("Copying slice")
// 	copy(s2, s1) // 将s1  的值拷贝到s2 中
// 	printSlice(s2)

// 	fmt.Println("Deleting slice")
// 	s2 = append(s2[0:3], s2[4:]...)
// 	printSlice(s2)
// }

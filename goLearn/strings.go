package main

// func main() {
// 	s := "Yes我爱慕课网！"
// 	// fmt.Println(len(s))

// 	// fmt.Println("Rune count:", utf8.RuneCountInString(s)) // 9

// 	bytes := []byte(s)
// 	for len(bytes) > 0 {
// 		_, size := utf8.DecodeRune(bytes)
// 		bytes = bytes[size:]
// 		// fmt.Printf("%c", r) // Yes我爱慕课网！
// 	}
// 	// fmt.Println()

// 	for i, ch := range []rune(s) {
// 		fmt.Printf("(%d %c)", i, ch)  // (0 Y)(1 e)(2 s)(3 我)(4 爱)(5 慕)(6 课)(7 网)(8 ！)
// 	}
// 	fmt.Println()
// }

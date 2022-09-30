package loop

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	printFileContents(file)
}

func forever() {
	for {
		fmt.Println("abc")
	}
}

func printFileContents(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func main() {
	fmt.Println()
	printFile("abc.txt")
	// ``跨行的文件
	s := `abc"d"
	kkkkkk
	
	123
	p`
	printFileContents(strings.NewReader(s))
}

package main

import (
	"fmt"

	"github.com/lithammer/shortuuid/v4"
)

func main() {
	id := shortuuid.New()

	fmt.Println("id:", id, "length:", len(id))

	id = shortuuid.NewWithNamespace("http://127.0.0.1.com")

	fmt.Println("id:", id, "length:", len(id))

	str := "12345#$%^&*67890qwerty/;'~!@uiopasdfghjklzxcvbnm,.()_+·><"

	id = shortuuid.NewWithAlphabet(str)
	fmt.Println("id:", id, "length:", len(id))

}

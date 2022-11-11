package main

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

func main() {
	id := ksuid.New()

	fmt.Println("id:", id, "length:", len(id))

	id1 := ksuid.New()
	id2 := ksuid.New()

	fmt.Println(id1, id2)

	compareResult := ksuid.Compare(id1, id2)
	fmt.Println(compareResult)

	isSorted := ksuid.IsSorted([]ksuid.KSUID{id2, id1})
	fmt.Println(isSorted)
}

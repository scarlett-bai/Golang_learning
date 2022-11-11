package main

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func main() {
	id, err := uuid.NewV1()
	if err != nil {
		fmt.Printf("uuid NewUUID err:%+v", err)
	}
	fmt.Println("id:", id.String(), "length:", len(id.String()))

	id, err = uuid.NewV4()
	if err != nil {
		fmt.Printf("uuid NewUUID err:%+v", err)
	}
	fmt.Println("id:", id.String(), "length:", len(id.String()))
}

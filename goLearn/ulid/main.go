package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func main() {
	t := time.Now().UTC()
	fmt.Println("time:", t)

	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	fmt.Println("id:", id.String(), "length:", len(id.String()))
}

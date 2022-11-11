package main

import (
	"crypto/md5"
	"fmt"

	"github.com/rs/xid"
)

func main() {
	id := xid.New()
	containerName := "test"
	containerNameID := make([]byte, 3)
	hw := md5.New()
	hw.Write([]byte(containerName))
	copy(containerNameID, hw.Sum(nil))
	id[4] = containerNameID[0]
	id[5] = containerNameID[1]
	id[6] = containerNameID[2]

	fmt.Println("id:", id, "length:", len(id))

}

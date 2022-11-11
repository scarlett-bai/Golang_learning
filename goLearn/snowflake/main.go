package main

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

func main() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	id := node.Generate().String()
	fmt.Println("id:", id, "length:", len(id))
}

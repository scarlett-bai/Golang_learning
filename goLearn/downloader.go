package main

import (
	"fmt"
	"goLearn/infra"
	"io/ioutil"
	"net/http"
)

func retrieve(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("%s\n", bytes)
	return bytes
}

func getRetrieve() retriver {
	return infra.Retriever{}
}

type retriver interface {
	Get(string) string
}

func main() {
	// bytes := retrieve("https://www.imooc.com")
	// fmt.Printf("%s\n", bytes)
	// retriever := getRetrieve()
	// var retriever infra.Retriever = getRetrieve()
	var r retriver = getRetrieve()

	fmt.Println(r.Get("https://www.imooc.com"))
}

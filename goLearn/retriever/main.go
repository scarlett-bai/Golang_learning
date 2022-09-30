package main

import (
	"fmt"
	"goLearn/retriever/mock"
	"goLearn/retriever/real"
)

type Retriever interface {
	Get(url string) string
}

type Poster interface {
	Post(url string, form map[string]string) string
}

type RetrieverPoster interface {
	Retriever
	Poster
}

const url = "https://www.imooc.com"

func session(s RetrieverPoster) string {
	s.Post(url, map[string]string{
		"contents": "another faked imooc.com",
	})
	return s.Get(url)
}

func main() {
	// var w io.ReadWriteCloser
	var r Retriever
	retriever := mock.Retriever{"this is a fake imooc.com"}
	r = &retriever
	inspect(r)
	fmt.Println(session(&retriever))
}

func inspect(r Retriever) {
	fmt.Println("Inspecting", r)
	fmt.Printf("> %T, %+v\n", r, r)
	fmt.Print("> Type switch:")
	switch v := r.(type) {
	case *mock.Retriever:
		fmt.Println("Contents：", v.Contents)
	case *real.Retriever:
		fmt.Println("UserAgent:", v.UserAgent)
	}
	fmt.Println()
}

package main

import (
	"goLearn/errhandling/filelistingserver/filelisting"
	"log"
	"net/http"
	"os"
)

type appHandler func(writer http.ResponseWriter, req *http.Request) error

func errWrapper(handler appHandler) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		// 处理panic
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic: %v\n", r)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

		}()
		err := handler(writer, request)
		if err != nil {
			log.Printf("Error occurred handling request:%s\n", err.Error())
			if userErr, ok := err.(userError); ok {
				//  处理 user error  可以把具体错误返回给用户
				http.Error(writer, userErr.Message(), http.StatusBadRequest)
				return
			}
			code := http.StatusOK
			// 系统Error  这个是不能直接返回给用户的
			switch {
			case os.IsNotExist(err):
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}
			http.Error(writer, http.StatusText(code), code)
		}
	}
}

type userError interface {
	error
	Message() string
}

func main() {
	http.HandleFunc("/", errWrapper(filelisting.HanleFileList))
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		panic(err)
	}
}

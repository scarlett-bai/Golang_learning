package main

import (
	"context"
	trippb "coolcar/proto/gen/go"
	trip "coolcar/tripservice"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.Lshortfile)
	go startGRPCGateway()
	log.Println("start")
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("finish lis")
	s := grpc.NewServer()
	trippb.RegisterTripServiceServer(s, &trip.Service{})
	log.Println("finish grpc")
	log.Fatal(s.Serve(lis))
	log.Println("GRPC start to serving")
}

func startGRPCGateway() {

	c := context.Background()          // 生成一个没有什么具体内容的上下文
	c, cancel := context.WithCancel(c) // 这么一个context 上下文 有 cancel 功能
	defer cancel()                     // 只要调了cancel() 就算传输了一半 也会被cancel掉

	// 因为mux 以后还会需要使用，所以就把它提出来
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			EnumsAsInts: true,
			OrigName:    true,
		},
	))
	err := trippb.RegisterTripServiceHandlerFromEndpoint(
		c, // 通过这个context 去调连接
		//mux: multiplexer 一对多 就是个分发器
		mux,                                    // 连接注册在 runtime.NewServeMux()里
		":8081",                                // 连接的地址
		[]grpc.DialOption{grpc.WithInsecure()}, // 连接的方式
	)
	if err != nil {
		log.Fatalf("cannot start grpc gateway: %v\n", err)
	}

	// 开始http监听
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("cannont listen and server: %v\n", err)
	}

}

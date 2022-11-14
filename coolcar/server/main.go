package main

import (
	trippb "coolcar/proto/gen/go"
	trip "coolcar/tripservice"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
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

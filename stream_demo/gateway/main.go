package main

import (
	"google.golang.org/grpc"
	"grpc-demo/stream_demo/pb"
	"grpc-demo/stream_demo/server"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听失败:%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, server.NewStreamServer())
	log.Println("监听端口50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动错误:%v", err)
	}
}

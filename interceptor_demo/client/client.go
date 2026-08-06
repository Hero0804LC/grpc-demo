package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	pb "grpc-demo/interceptor_demo/pb/proto"
	"grpc-demo/interceptor_demo/server"
	"log"
	"time"
)

func main() {
	ready := make(chan struct{})
	go func() {
		server.Start(ready)
	}()
	<-ready
	conn, err := grpc.NewClient("localhost:50061", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewHelloServiceClient(conn)
	//注入token
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer valid-token")
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{
		Name: "world",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("response: %s", resp.Message)
}

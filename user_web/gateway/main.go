package main

import (
	"context"
	"google.golang.org/grpc/status"
	"grpc-demo/user_web/server"
	"log"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-demo/user_web/pb"
)

func main() {
	ready := make(chan struct{})
	// 启动 gRPC Server
	go func() {
		server.StartGRPCServer(ready)
	}()
	<-ready
	// 等 gRPC 启动
	time.Sleep(time.Second)

	r := gin.Default()

	// gRPC Client
	addr := net.JoinHostPort("127.0.0.1", "50051")
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// HTTP -> gRPC
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.GetInt64("id")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := client.GetUserById(ctx, &pb.IdRequest{Id: id})
		if err != nil {
			st, _ := status.FromError(err)
			c.JSON(404, gin.H{"error": st.Message()})
			return
		}

		c.JSON(200, resp)
	})

	r.GET("/users", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := client.GetUserList(ctx, &pb.PageInfo{
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, resp)
	})

	log.Println("Gin gateway listening on :8080")
	r.Run(":8080")
}

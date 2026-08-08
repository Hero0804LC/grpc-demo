package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "grpc-demo/seckill/pb"
)

func main() {
	conn, err := grpc.NewClient(
		"127.0.0.1:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer conn.Close()

	client := pb.NewSeckillServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.DoSeckill(ctx, &pb.SeckillRequest{
		UserId:  1,
		GoodsId: 1001,
	})
	if err != nil {
		log.Fatalf("抢购失败: %v", err)
	}

	log.Printf("抢购成功: orderNo=%s msg=%s", resp.OrderNo, resp.Msg)
}

package main

import (
	"context"
	pb2 "grpc-demo/stream_demo/pb"
	pb1 "grpc-demo/user_web/pb/proto"
	"io"
	"log"
	"time"
)

// stream_server
func main() {
	//获取全局grpc客户端
	c := Client()
	conn, err := c.Conn()
	if err != nil {
		log.Fatalf("grpc连接失败:%v", err)
	}
	defer c.Close()
	streamClient := pb2.NewStreamServiceClient(conn)
	ctx, cannel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cannel()
	stream, err := streamClient.SubscribeLogs(ctx, &pb2.LogSubscribeRequest{
		Topic: "access-log",
	})
	if err != nil {
		log.Fatalf("请求失败:%v", err)
	}
	for {
		entry, err := stream.Recv()
		if err == io.EOF {
			log.Println("stream closed by server")
			break
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}

		log.Printf("recv log: seq=%d msg=%s", entry.Seq, entry.Msg)
	}
}

// user_web
func main1() {
	//获取全局grpc客户端
	c := Client()
	conn, err := c.Conn()
	if err != nil {
		log.Fatalf("grpc连接失败:%v", err)
	}
	defer c.Close()

	//创建用户客户端
	userClient := pb1.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	/*修改用户*/
	getResp, err := userClient.UpdateUser(ctx, &pb1.UpdateUserInfo{
		Id:       1,
		NickName: "修改名字",
	})
	if err != nil {
		log.Fatalf("修改用户信息出错%v", err)
	}
	log.Printf("修改用户信息成功%v", getResp)
	/*通过手机号获取用户
	getResp, err := userClient.GetUserByMobile(ctx, &pb1.MobileRequest{
		Mobile: "15066666666",
	})
	if err != nil {
		log.Fatalf("通过手机号获取用户失败%v", err)
	}
	log.Printf("通过手机号获取用户成功%v", getResp)*/
	/*通过ID获取用户信息
	getResp, err := userClient.GetUserById(ctx, &pb1.IdRequest{
		Id: 1,
	})
	if err != nil {
		log.Fatalf("通过ID获取用户失败:%v", err)
	}
	log.Printf("通过ID获取用户成功%+v", getResp)*/
	/*创建用户
	createResp, err := userClient.CreateUser(ctx, &pb1.CreateUserInfo{
		NickName: "张三",
		Password: "123456",
		Mobile:   "15066666666",
	})
	if err != nil {
		log.Fatalf("创建用户失败:%v", err)
	}
	log.Printf("创建用户成功%v", createResp)*/
}

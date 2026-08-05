package main

import (
	"grpc-demo/user_web/server"
	"log"
)

func main() {
	log.Println("启动服务")
	server.StartGRPCServer()
}

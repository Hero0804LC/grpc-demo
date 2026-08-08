package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"grpc-demo/seckill/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ready := make(chan struct{})
	server.Start(ready)
	<-ready
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "seckill_order",
		GroupID: "order-consumer-group",
	})
	log.Println("开始消费")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-sig:
			log.Println("结束")
			r.Close()
			return
		default:
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				continue
			}
			orderNo, userId, goodsId := parse(string(m.Value))
			log.Printf("模拟写入数据库: %s, user=%d, goods=%d",
				orderNo, userId, goodsId)
		}
	}
}

// 解析
func parse(s string) (orderNo string, userId, goodsId int64) {
	fmt.Sscanf(s, "%[^,],%d,%d", &orderNo, &userId, &goodsId)
	return
}

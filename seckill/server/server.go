package server

import (
	"context"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-demo/seckill/pb"
	"log"
	"net"
	"strconv"
	"time"
)

type SeckillServer struct {
	pb.UnimplementedSeckillServiceServer
	kafkaWriter *kafka.Writer
}

func NewSeckillServer() *SeckillServer {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "seckill_order",
		Balancer: &kafka.LeastBytes{},
	}
	return &SeckillServer{kafkaWriter: w}
}

// 内存库存
var stock = map[int64]int64{
	1001: 100,
}

// 执行抢购
func (s *SeckillServer) DoSeckill(ctx context.Context, req *pb.SeckillRequest) (*pb.SeckillResponse, error) {
	//校验库存
	if stock[req.GoodsId] <= 0 {
		return &pb.SeckillResponse{
			Msg: "没有库存",
		}, nil
	}
	stock[req.GoodsId]--
	//生成订单号
	orderNo := strconv.FormatInt(time.Now().UnixNano(), 10)
	//生产信息
	msg := []byte(orderNo + "," + strconv.FormatInt(req.UserId, 10) + "," + strconv.FormatInt(req.GoodsId, 10))
	err := s.kafkaWriter.WriteMessages(ctx, kafka.Message{Value: msg})
	if err != nil {
		log.Println("kafka 写入失败:", err)
		return nil, err
	}

	return &pb.SeckillResponse{
		OrderNo: orderNo,
		Msg:     "抢购成功",
	}, nil
}
func Start(ready chan<- struct{}) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSeckillServiceServer(s, NewSeckillServer())
	reflection.Register(s)

	go func() {
		log.Println("seckill gRPC server listening on :50051")
		close(ready)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("grpc server error: %v", err)
		}
	}()
}

package server

import (
	"grpc-demo/stream_demo/pb"
	"log"
	"strconv"
	"time"
)

type StreamServer struct {
	pb.UnimplementedStreamServiceServer
}

func NewStreamServer() *StreamServer {
	return &StreamServer{}
}

func (s *StreamServer) SubscribeLogs(req *pb.LogSubscribeRequest, stream pb.StreamService_SubscribeLogsServer) error {
	log.Printf("client subscribed topic: %s", req.Topic)
	ctx := stream.Context()
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			log.Println("client disconnected")
		default:
			entry := &pb.LogEntry{
				Seq:   int64(i),
				Level: "INFO",
				Msg:   "log line" + strconv.Itoa(i),
				Ts:    time.Now().Unix(),
			}
			if err := stream.Send(entry); err != nil {
				return err
			}
			log.Printf("sent log seq=%d", i)
			time.Sleep(time.Second)
		}
	}
	log.Println("log stream finished")
	return nil
}

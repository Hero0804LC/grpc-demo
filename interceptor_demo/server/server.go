package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "grpc-demo/interceptor_demo/pb/proto"
	"log"
	"net"
)

type helloServer struct {
	pb.UnimplementedHelloServiceServer
}

func (s *helloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello" + req.Name,
	}, nil
}

func Start(ready chan<- struct{}) {
	lis, err := net.Listen("tcp", ":50061")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(recoveryInterceptor, loggingInterceptor, authInterceptor))
	pb.RegisterHelloServiceServer(s, &helloServer{})
	reflection.Register(s)
	go func() {
		log.Println("interceptor-demo server listening on :50061")
		close(ready)
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

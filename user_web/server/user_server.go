package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "grpc-demo/user_web/pb"
	"log"
	"net"
	"sync"
)

// 模拟数据库user
type userModel struct {
	ID       int64
	NickName string
	Password string
	Mobile   string
	Gender   string
	Birthday *timestamppb.Timestamp
	Role     int64
}

var (
	users        = make(map[int64]*userModel)
	nextID int64 = 1
	mu     sync.Mutex
)

// gRPC Server 实现

type userServer struct {
	pb.UnimplementedUserServiceServer
}

func StartGRPCServer(ready chan<- struct{}) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &userServer{})
	reflection.Register(s)

	go func() {
		log.Println("gRPC server listening on :50051")
		close(ready)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

// 获取用户列表
func (s *userServer) GetUserList(ctx context.Context, req *pb.PageInfo) (*pb.UserListResponse, error) {
	mu.Lock()
	defer mu.Unlock()
	var list []*pb.UserInfoResponse
	for _, u := range users {
		list = append(list, toPB(u))
	}
	return &pb.UserListResponse{
		Total: int64(len(list)),
		Data:  list,
	}, nil
}
func toPB(u *userModel) *pb.UserInfoResponse {
	return &pb.UserInfoResponse{
		Id:       u.ID,
		Mobile:   u.Mobile,
		NickName: u.NickName,
		Gender:   u.Gender,
		Birthday: u.Birthday,
		Role:     u.Role,
	}
}

// 通过ID获取用户
func (s *userServer) GetUserById(ctx context.Context, req *pb.IdRequest) (*pb.UserInfoResponse, error) {
	mu.Lock()
	defer mu.Unlock()
	u, ok := users[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "用户未找到")
	}
	return toPB(u), nil
}

// 通过手机号获取用户
func (s *userServer) GetUserByMobile(ctx context.Context, req *pb.MobileRequest) (*pb.UserInfoResponse, error) {
	mu.Lock()
	defer mu.Unlock()
	for _, u := range users {
		if u.Mobile == req.Mobile {
			return toPB(u), nil
		}
	}
	return nil, status.Error(codes.NotFound, "用户未找到")
}

// 创建用户
func (s *userServer) CreateUser(ctx context.Context, req *pb.CreateUserInfo) (*pb.CreateUserResponse, error) {
	mu.Lock()
	defer mu.Unlock()
	//数据非空校验
	if req.NickName == "" || req.Password == "" || req.Mobile == "" {
		return nil, status.Error(codes.InvalidArgument, "缺失注册信息")
	}
	id := nextID
	nextID++
	users[id] = &userModel{
		ID:       id,
		NickName: req.NickName,
		Password: "加密拼接" + req.Password,
		Mobile:   req.Mobile,
		Role:     1,
	}
	return &pb.CreateUserResponse{
		Id: id,
	}, nil
}

// 修改用户
func (s *userServer) UpdateUser(ctx context.Context, req *pb.UpdateUserInfo) (*pb.UpdateUserResponse, error) {
	mu.Lock()
	defer mu.Unlock()
	u, ok := users[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "未找到用户")
	}
	if req.NickName != "" {
		u.NickName = req.NickName
	}
	if req.Gender != "" {
		u.Gender = req.Gender
	}
	if req.Birthday != nil {
		u.Birthday = req.Birthday
	}
	return &pb.UpdateUserResponse{Success: true}, nil
}

// 删除用户
func (s *userServer) DeleteUser(ctx context.Context, req *pb.IdRequest) (*emptypb.Empty, error) {
	log.Printf("DeleteUser: id=%d", req.Id)
	mu.Lock()
	defer mu.Unlock()
	if _, ok := users[req.Id]; !ok {
		return nil, status.Error(codes.NotFound, "未找到用户")
	}
	delete(users, req.Id)
	return &emptypb.Empty{}, nil
}

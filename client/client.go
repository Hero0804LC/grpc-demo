package main

import (
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultTarget = "localhost:50051"

type GRPCClient struct {
	conn *grpc.ClientConn
	mu   sync.Mutex
}

var (
	globalClient *GRPCClient
	once         sync.Once
)

// Client 返回全局 gRPC 客户端（单例）
func Client() *GRPCClient {
	once.Do(func() {
		globalClient = &GRPCClient{}
	})
	return globalClient
}

// Conn 获取 gRPC 连接（懒加载）
func (c *GRPCClient) Conn(target ...string) (*grpc.ClientConn, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 已存在连接，直接返回
	if c.conn != nil {
		return c.conn, nil
	}

	addr := defaultTarget
	if len(target) > 0 && target[0] != "" {
		addr = target[0]
	}

	// 不立即建立 TCP 连接，首次 RPC 时自动连接
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	log.Printf("[gRPC Client] created client for %s", addr)
	return conn, nil
}

// Close 关闭连接（程序退出时调用）
func (c *GRPCClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		log.Println("[gRPC Client] connection closed")
		return err
	}
	return nil
}

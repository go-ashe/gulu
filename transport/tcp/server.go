package tcp

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/go-ashe/gulu/logx"
)

type TCPServer struct {
	conf     *Config
	listener net.Listener
	handler  func(net.Conn)
	mu       sync.RWMutex
}

// NewTCPServer creates a new TCP server.
func NewTCPServer(ctx context.Context, conf *Config) (*TCPServer, error) {
	return &TCPServer{
		conf: conf,
	}, nil
}

// SetHandler 设置连接处理函数
func (s *TCPServer) SetHandler(handler func(net.Conn)) {
	s.handler = handler
}

// Start starts the TCP server.
func (s *TCPServer) Start(_ context.Context) error {
	addr := fmt.Sprintf(":%d", s.conf.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("启动服务失败: %v", err)
	}

	s.listener = listener
	logx.L().Info("服务器启动成功", "listen_addr", addr)

	go func() {
		for {
			conn, err2 := s.listener.Accept()
			if err2 != nil {
				logx.L().Error("接受连接失败", "error", err2.Error())
				continue
			}

			go s.handleConnection(conn)
		}
	}()

	return nil
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// 处理连接
	s.handler(conn)

}

// Stop stops the TCP server.
func (s *TCPServer) Stop(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
	return nil
}

package udp

import (
	"context"
	"fmt"
	"log"
	"net"
)

type Handler func(*net.UDPAddr, []byte) []byte

// ServerConfig UDP服务器配置
type ServerConfig struct {
	// 监听端口
	Port int
	// 读取缓冲区大小
	BufferSize int
}

// Server UDP服务器
type Server struct {
	config  *ServerConfig
	conn    *net.UDPConn
	handler Handler
}

type ClientConfig struct {
	Address string `json:"address" yaml:"address"`
}

// NewServer 创建一个新的UDP服务器
func NewServer(config *ServerConfig, h Handler) (s *Server, err error) {
	s = &Server{
		config:  config,
		handler: h,
	}

	if config.BufferSize <= 0 {
		config.BufferSize = 4096
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: config.Port,
	})
	if err != nil {
		err = fmt.Errorf("监听地址失败: %v", err)
		return
	}

	s.conn = conn
	return
}

// Start 启动服务器
func (s *Server) Start(_ context.Context) error {
	log.Printf("UDP服务器启动成功，监听端口: %d", s.config.Port)
	go s.serve()
	return nil
}

// serve 处理UDP消息
func (s *Server) serve() {
	buffer := make([]byte, s.config.BufferSize)
	for {
		n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("读取数据失败: %v", err)
			continue
		}

		// 复制数据以避免被后续读取覆盖
		data := make([]byte, n)
		copy(data, buffer[:n])

		// 处理消息
		if s.handler != nil {
			go s.handleMessage(remoteAddr, data)
		}
	}
}

// handleMessage 处理单个UDP消息
func (s *Server) handleMessage(addr *net.UDPAddr, msg []byte) {
	if response := s.handler(addr, msg); response != nil {
		if _, err := s.conn.WriteToUDP(response, addr); err != nil {
			log.Printf("发送响应失败: %v", err)
		}
	}
}

// Stop 关闭服务器
func (s *Server) Stop(_ context.Context) error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// LocalAddr 返回服务器监听的本地地址
func (s *Server) LocalAddr() net.Addr {
	if s.conn == nil {
		return nil
	}
	return s.conn.LocalAddr()
}

func (s *Server) GetConn() *net.UDPConn {
	return s.conn
}

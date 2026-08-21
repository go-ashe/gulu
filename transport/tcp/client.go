package tcp

import (
	"context"
	"fmt"
	"net"
	"sync"
)

// TCPClient represents a TCP client.
type TCPClient struct {
	conn net.Conn
	mu   sync.Mutex
}

// NewTCPClient creates a new TCP client and connects to the server.
func NewTCPClient(addr string) (*TCPClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &TCPClient{
		conn: conn,
	}, nil
}

// Send sends data to the server.
func (c *TCPClient) Send(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.conn.Write(data)
	return err
}

// Receive receives data from the server.
func (c *TCPClient) Receive(buffer []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.Read(buffer)
}

// SendAndReceive sends data to the server and receives data from the server.
func (c *TCPClient) SendAndReceive(_ context.Context, sendData []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 发送数据
	_, err := c.conn.Write(sendData)
	if err != nil {
		return nil, fmt.Errorf("发送数据失败: %v", err)
	}

	// 使用动态缓冲区读取响应
	buffer := make([]byte, 1024) // 初始缓冲区
	var response []byte
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			if len(response) > 0 {
				break // 已经读取到一些数据，可以返回
			}
			return nil, fmt.Errorf("读取响应失败: %v", err)
		}
		response = append(response, buffer[:n]...)
		if n < len(buffer) {
			break // 本次读取的数据小于缓冲区，说明数据已读完
		}
	}

	return response, nil
}

// Close closes the TCP connection.
func (c *TCPClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

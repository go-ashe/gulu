package udp

import (
	"fmt"
	"net"
)

type UDPClient struct {
	conn *net.UDPConn
}

// NewClient 创建一个新的 UDP 客户端
func NewClient(conn *net.UDPConn) *UDPClient {
	return &UDPClient{
		conn: conn,
	}
}

// SendTo 发送数据到指定地址
func (s *UDPClient) SendTo(data []byte, addr string) error {
	defer func() {
		//log.Printf("send to remote agv, addr:%s, data:%s", addr, data)
	}()

	if s.conn == nil {
		return fmt.Errorf("client conn is nil")
	}

	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("解析目标地址失败: %v", err)
	}

	if _, err = s.conn.WriteToUDP(data, raddr); err != nil {
		return fmt.Errorf("发送数据失败: %v", err)
	}

	return nil
}

// Receive 接收数据，支持超时
func (c *UDPClient) Receive() ([]byte, *net.UDPAddr, error) {
	// TODO: 缓冲区大小
	buffer := make([]byte, 1024)

	n, addr, err := c.conn.ReadFromUDP(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, nil, fmt.Errorf("读取超时")
		}
		return nil, nil, fmt.Errorf("读取数据失败: %v", err)
	}

	// 创建新的切片以返回数据，因为原缓冲区会被放回池中
	data := make([]byte, n)
	copy(data, buffer[:n])

	return data, addr, nil
}

// Close 关闭客户端
func (c *UDPClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetLocalAddr 获取本地地址
func (c *UDPClient) GetLocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

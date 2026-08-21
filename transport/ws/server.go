package ws

import (
	"github.com/go-ashe/gulu/slogx"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Server struct {
	Upgrader *websocket.Upgrader
	clients  *sync.Map
	Handler  func(c *Client) error
}

func New(h func(c *Client) error) *Server {
	return &Server{
		Upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源
			},
		},
		clients: new(sync.Map),
		Handler: h,
	}
}

// CloseClient 关闭单个客户端连接
func (s *Server) CloseClient(cli *Client) {
	if cli.Conn == nil {
		// 从 clients 中移除该连接
		s.clients.Delete(cli.Id)
		return
	}

	// 发送关闭消息并关闭连接
	if err := cli.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		// 如果连接已经关闭，忽略 "websocket: already closed" 错误
		if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			slogx.L().Error("[WebSocket] Error sending close message", "error", err)
		}
	}

	if err := cli.Conn.Close(); err != nil {
		slogx.L().Error("[WebSocket] Error closing WebSocket client", "error", err.Error())
	}

	// 从 clients 中移除该连接
	s.clients.Delete(cli.Id)
}

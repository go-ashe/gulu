package ws

import (
	"fmt"
	"github.com/go-ashe/gulu/slogx"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Client struct {
	Id            string
	Conn          *websocket.Conn
	MsgIdentifier string
	SendCh        chan *Message
	CloseOnce     sync.Once
}

func (s *Server) AddClient(cli *Client) {
	s.clients.Store(cli.Id, cli)
}

func (c *Client) WritePump(s *Server) {
	ticker := time.NewTicker(15 * time.Second)
	defer func() {
		ticker.Stop()
		c.CloseOnce.Do(func() {
			c.Close(s)
		})
	}()

	for {
		select {
		case msg, ok := <-c.SendCh:
			if !ok {
				return
			}
			if err := c.Conn.WriteJSON(msg); err != nil {
				slogx.L().Error("[WebSocket] Error writing message", "error", err.Error())
				return
			}
		case <-ticker.C:
			// ping 心跳
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slogx.L().Error("[WebSocket] WritePump ping error", "error", err.Error())
				return
			}
		}
	}
}

// ReadPump 读取消息，主要用来处理客户端主动关闭
func (c *Client) ReadPump(s *Server) {
	var err error
	c.Conn.SetReadLimit(512)
	if err = c.Conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		slogx.L().Error("[WebSocket] ReadPump SetReadDeadline error", "error", err.Error())
	}
	c.Conn.SetPongHandler(func(string) error {
		if err = c.Conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
			slogx.L().Error("[WebSocket] ReadPump SetPongHandler SetReadDeadline error", "error", err.Error())
		}
		return nil
	})

	for {
		if s.Handler != nil {
			if err = s.Handler(c); err != nil {
				if !websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseUnsupportedData,
					websocket.CloseNoStatusReceived) {
					slogx.L().Error("[WebSocket] ReadPump client disconnected", "id", c.Id)
				}
				c.Close(s)
				break
			}
		}
	}
}

func (c *Client) Close(s *Server) {
	c.CloseOnce.Do(func() {
		s.clients.Delete(c.Id)
		close(c.SendCh)
		err := c.Conn.Close()
		slogx.L().Info(fmt.Sprintf("[WebSocket] Client[%s] %s closed, error:%v", c.Conn.RemoteAddr().String(), c.Id, err))
	})
}

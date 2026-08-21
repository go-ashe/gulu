package ws

import (
	"github.com/pkg/errors"
	"log"
)

type Message struct {
	Identifier  string      `json:"-"` // 消息唯一标识
	MessageType int         `json:"messageType"`
	Code        int         `json:"code"`
	Msg         string      `json:"msg"`
	Data        interface{} `json:"data"`
}

func (s *Server) SendBroadcastMsg(msg *Message) error {
	if msg.MessageType == 0 || len(msg.Identifier) == 0 {
		return errors.New("[WebSocket] message type or identifier is empty")
	}
	s.clients.Range(func(key, val interface{}) bool {
		cli, ok := val.(*Client)
		if !ok || cli.MsgIdentifier != msg.Identifier {
			return true
		}

		if cli.Conn == nil {
			s.clients.Delete(cli.Id)
			return true
		}

		if err := cli.Conn.WriteJSON(msg); err != nil {
			log.Printf("[WebSocket] Error writing message: %v", err)
			s.CloseClient(cli)
		}
		return true
	})

	return nil
}

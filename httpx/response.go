package httpx

import (
	"fmt"
	"reflect"
)

type Json interface {
	OK() Json
	WithCode(code int) Json
	WithData(data interface{}) Json
	WithMessageType(msgType int) Json
	Error(err error) Json
}

type JsonRsp struct {
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	MessageType int    `json:"messageType,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Data        any    `json:"data,omitempty"`
}

func NewJson() Json {
	return &JsonRsp{}
}

// JsonResponse 输出Json数据
//
//	args 应该是 [data], [data,err]或者空
//	err 必须是 error
func JsonResponse(args ...any) Json {
	rsp := NewJson()
	if len(args) == 1 {
		if args[0] != nil {
			return rsp.OK().WithData(args[0])
		}
		return rsp.OK()
	}
	if len(args) == 2 {
		if args[0] != nil {
			rsp.WithData(args[0])
		}
		if args[1] != nil {
			if e, ok := args[1].(error); ok {
				return rsp.Error(e)
			}
		}
	}
	return rsp.OK()
}

func (r *JsonRsp) OK() Json {
	r.Code = 0
	r.Msg = "OK"
	return r
}

func (r *JsonRsp) WithCode(code int) Json {
	r.Code = code
	return r
}

func (r *JsonRsp) WithData(data any) Json {
	if reflect.ValueOf(data).IsZero() {
		return r
	}
	r.Data = data
	return r
}

func (r *JsonRsp) WithMessageType(messageType int) Json {
	r.MessageType = messageType
	return r
}

func (r *JsonRsp) Error(err error) Json {
	e := FromError(err)
	if e == nil {
		r.Code = 500
		r.Msg = ""
		return r
	}

	r.Code = e.Code
	r.Msg = e.Message
	r.Reason = e.Reason
	if len(e.Metadata) > 0 {
		r.Reason += fmt.Sprintf("| metadata:%v", e.Metadata)
	}

	return r
}

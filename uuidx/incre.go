package uuidx

import (
	"sync/atomic"
	"time"
)

var (
	uint32SeqPool = make(map[string]*uint32, 100)
)

func GetUint32Seq(uniKey string) uint32 {
	if uint32SeqPool[uniKey] == nil {
		uint32SeqPool[uniKey] = new(uint32)
		atomic.AddUint32(uint32SeqPool[uniKey], uint32(time.Now().Unix()%100000))
	}
	return atomic.AddUint32(uint32SeqPool[uniKey], 1)
}

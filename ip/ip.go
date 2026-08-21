package ip

import (
	"fmt"
	"math/rand"
	"net"
)

func GetRandomIPV4() string {
	// 生成随机ip:port
	return fmt.Sprintf("10.10.%d.%d:%d", rand.Intn(255), rand.Intn(255), int64(10000+rand.Intn(10000)))
}

func GetLocalMachineIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

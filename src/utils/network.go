//Package utils 文件服务器
package utils

import "net"

//IPV4 get V4
func IPV4() string {

	result := ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = ipnet.IP.String()
				break
			}
		}
	}
	return result
}

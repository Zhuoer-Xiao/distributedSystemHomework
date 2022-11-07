package common

import (
	"fmt"
	"net"
	"strings"
)

//分配端口
const(
	HeartBeatPort =iota+626
)

//获取本机IP
func LocalIp()net.IP{
	conn,err:=net.Dial("udp","www.baidu.com:80")
	if(err!=nil){
		fmt.Println("local Ip error:",err)
	}
	defer conn.Close()
	return net.ParseIP(strings.Split(conn.LocalAddr().String(),":")[0])
}


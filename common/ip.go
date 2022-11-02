package common

import (
	"fmt"
	"net"
	"strings"
)

func LocalIp()net.IP{
	conn,err:=net.Dial("udp","www.baidu.com:80")
	if(err!=nil){
		fmt.Println("local Ip error:",err)
	}
	defer conn.Close()
	return net.ParseIP(strings.Split(conn.LocalAddr().String(),":")[0])
}


package common

import (
	"fmt"
	"net"
	"net/rpc"
	"strings"
)

//分配端口
const (
	HeartBeatPort =iota+1233
	ManagerPort
)

var (
	MasterIp=net.ParseIP("127.0.0.1")
	Conn  *rpc.Client
)

func init(){
	addr := fmt.Sprintf("%s:%v", MasterIp.String(), ManagerPort)
	Conn, _ = rpc.Dial("tcp", addr)
}

//获取本机IP
func LocalIp()net.IP{
	conn,err:=net.Dial("udp","www.baidu.com:80")
	if(err!=nil){
		fmt.Println("local Ip error:",err)
	}
	defer conn.Close()
	return net.ParseIP(strings.Split(conn.LocalAddr().String(),":")[0])
}




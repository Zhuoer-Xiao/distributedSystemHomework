package common
import(
	"net"
)
//此处定义各种rpc消息体
type HeartBeatArgs struct{
	IP net.IP
	Port int
}
type HeartBeatReply struct{
	
}
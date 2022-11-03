package common

import (
	"net"
)

// 此处定义各种rpc消息体
type HeartBeatArgs struct {
	IP   net.IP
	Port int
}
type HeartBeatReply struct {
}

// 打开文件
type OpenArgs struct {
	FileName string //文件名
	Index    int    //所在chunk块
	Perm     uint32 //权限
}

type OpenReply struct {
	chunkName uint64
}

// 关闭文件
type CloseArgs struct {
}

type CloseReply struct {
}

// 读取文件
type ReadArgs struct {
}

type ReadReply struct {
}

// delete
type DeleteArgs struct {
}

type DeleteReply struct {
}

// create
type CreateArgs struct {
}

type CreateReply struct {
}

// append
type AppendArgs struct {
}

type AppendReply struct {
}

// exist
type ExistArgs struct {
}

type ExistReply struct {
}

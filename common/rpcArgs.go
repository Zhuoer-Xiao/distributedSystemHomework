package common

import (
	"distributedSystemHomework/chunkserver"
	"net"
	//"distributedSystemHomework/chunkserver"
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

// 待修改
type OpenReply struct {
	ChunkName         uint64 //块id
	ChunkServerNameIp string //块所在chunkserver信息
}

// 更新元数据
type UpdateArgs struct {
	Files []*File
}

type UpdateReply struct {
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
	Test1 string
	Test2 string
}

type CreateReply struct {
	TestRes string
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

type CreateChunkRpcArgs struct {
	Path string
}

type CreateChunkRpcReply struct {
	Addresses []ServerAddress
	Handle ChunkHandle
}

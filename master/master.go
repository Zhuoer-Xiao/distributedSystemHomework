package master

import (
	// "errors"
	// "fmt"
	// "log"
	// "net/rpc"
	"distributedSystemHomework/chunkserver"
)

//命名空间
//整个文件系统的目录结构和Chunk的基本信息；
//文件与Chunk的映射关系；
//各个Chunk备份（Replicas）的位置信息，默认为3个副本。
//根据文件与偏移找到所在chunk

// 心跳机制(可不考虑)
//心跳机制是让Master服务器了解到Chunk服务器的状态，检测Chunk服务器是否在线以及获取相关信息。

//故障转移(可不考虑)
//

// master的数据结构
type Master struct {
	chunkServers map[uint64]*chunkserver.ChunkServer //保存所有chunkserver信息
}

// 初始化master
func NewMaster() *Master {
	m := new(Master)
	return m
}

func FindFileChunk(fileName string, fileIndex uint16) uint64 {
	return 0
}

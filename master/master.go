package master

import (
	// "errors"
	// "fmt"
	// "log"
	// "net/rpc"
	"distributedSystemHomework/chunkserver"
	"distributedSystemHomework/common"
	"fmt"
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
	chunkServers map[string]*chunkserver.ChunkServer //保存所有chunkserver信息，通过ip来标志chunkserver
	nameSpace    *NameSpace                          //命名空间
}

// 初始化master
func NewMaster() *Master {
	m := new(Master)
	m.chunkServers = make(map[string]*chunkserver.ChunkServer)
	m.nameSpace = NewNameSpace()
	return m
}

func (m *Master) OpenFile(args *common.OpenArgs, reply *common.OpenReply) error {
	if common.CheckCreate(args.Perm) {
		file, err := m.nameSpace.createFile(args.FileName, args.Index, args.Perm)
		if err != nil {
			fmt.Println("Open File: ", args.FileName, " fail.")
			return err
		}
	}
	return nil
}

//定期写入内存
//为chunkserver分配ip

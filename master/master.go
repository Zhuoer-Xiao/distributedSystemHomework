package master

import (
	// "errors"
	// "fmt"
	"log"
	"net/rpc"
	"distributedSystemHomework/chunkserver"
	"distributedSystemHomework/common"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"net"
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
// 待测试
type Master struct {
	chunkServers map[string]*chunkserver.ChunkServer //保存所有chunkserver信息，通过ip来标志chunkserver
	nameSpace    *NameSpace                          //命名空间
	openFiles []*common.File//已打开文件
}

// 初始化master
// 待测试
func NewMaster() *Master {
	m := new(Master)
	m.chunkServers = make(map[string]*chunkserver.ChunkServer)
	m.nameSpace = NewNameSpace()
	return m
}

// 找到client需要的文件chunk位置
// 待测试
func (m *Master) OpenFile(args *common.OpenArgs, reply *common.OpenReply) error {
	if common.CheckCreate(args.Perm) {
		file, err := m.nameSpace.CreateFile(args.FileName)
		if err != nil {
			fmt.Println("Open File: ", args.FileName, " fail.")
			return err
		}
		reply.ChunkName = file.Chunks[args.Index]
	}

	return nil
}

// 定期写入内存,转换为json写入
// 写入读出结构体
// 待测试
func (m *Master) WriteFileInfo() error {
	file1, _ := os.OpenFile("./file1.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	file2, _ := os.OpenFile("./file2.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	outPut1, _ := json.Marshal(&m.chunkServers)
	outPut2, _ := json.Marshal(&m.nameSpace)
	file1.Write(outPut1)
	file2.Write(outPut2)
	file1.Close()
	file2.Close()
	return nil
}

// 开启master处理
// 每隔十分钟将master数据写入内存
func (m *Master) Main() error {
	m.openHeartbeatServer()
	for i := 1; i <= 10; i++ {
		m.WriteFileInfo()
		time.Sleep(time.Minute*10)
	}
	return nil
}

//选择Chunkserver返回
//待完成
func(m* Master)PickChunkServer(chunk uint64) *chunkserver.ChunkServer{
	for _, cs := range m.chunkServers {
			return cs
	}

	return nil
}

//监听rpc信息
//待测试
func (m *Master) openHeartbeatServer() {
	r := rpc.NewServer()
	r.Register(m)

	addr := fmt.Sprintf(":%v", common.HeartBeatPort)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	go r.Accept(l)
}

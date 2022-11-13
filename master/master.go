package master

import (
	// "errors"
	// "fmt"
	"distributedSystemHomework/chunkserver"
	"distributedSystemHomework/common"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"time"
)

//命名空间
//整个文件系统的目录结构和Chunk的基本信息；
//文件与Chunk的映射关系；
//各个Chunk备份（Replicas）的位置信息，默认为3个副本。
//根据文件与偏移找到所在chunk

// 心跳机制
//心跳机制是让Master服务器了解到Chunk服务器的状态，检测Chunk服务器是否在线以及获取相关信息。
//此处心跳交换元数据信息

//故障转移(可不考虑)

//更新元数据
//待完成

// master的数据结构
// 待测试
type Master struct {
	chunkServers map[string]*chunkserver.ChunkServer //保存所有chunkserver信息，通过ip来标志chunkserver
	nameSpace    *NameSpace                          //命名空间

	chunksLocation map[uint64]*[3]string //chunk所对应的chunk文件
	uuid           *common.UuidGenerate
}

// 初始化master
// 待测试
func NewMaster() *Master {
	m := new(Master)
	m.chunkServers = make(map[string]*chunkserver.ChunkServer)
	m.nameSpace = NewNameSpace()
	m.chunksLocation = make(map[uint64]*[3]string)
	m.uuid = common.NewUuidGnerate()
	return m
}

// 找到client需要的文件chunk位置
// 待修改：index越界，权限问题,返回chunkserver信息
func (m *Master) OpenFile(args *common.GetChunkHandleArg, reply *common.GetChunkHandleReply) error {

	file, err := m.nameSpace.FindFile(string(args.Path))
	if err != nil {
		fmt.Println("Open File: ", args.Path, " fail.")
		return err
	}
	if len(file.Chunks) < int(args.Index) {
		return errors.New("Out of chunk's index")
	} else {
		reply.Handle = common.ChunkHandle(file.Chunks[args.Index])
		//reply.ChunkServerNameIp = *m.PickChunkServer(reply.ChunkName).address
	}

	//无权限
	return errors.New("No Permission")
}

// 定期写入内存,转换为json写入
// 写入读出结构体
// 已测试
func (m *Master) WriteFileInfo() error {
	file1, _ := os.OpenFile("./file1.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	file2, _ := os.OpenFile("./file2.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	file3, _ := os.OpenFile("./file3.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	outPut1, _ := json.Marshal(&m.chunkServers)
	outPut2, _ := json.Marshal(&m.nameSpace)
	outPut3, _ := json.Marshal(&m.chunksLocation)
	file1.Write(outPut1)
	file2.Write(outPut2)
	file3.Write(outPut3)
	file1.Close()
	file2.Close()
	file3.Close()
	return nil
}

// 开启master处理
// 每隔十分钟将master数据写入内存
// 已测试
func (m *Master) Main() error {
	m.openHeartbeatServer()
	for true {
		m.WriteFileInfo()
		time.Sleep(time.Minute * 10)
	}
	return nil
}

// 选择Chunkserver返回
// 待测试
// 默认只有三个副本，那么生成一个随机数来挑选chunkserver
func (m *Master) PickChunkServer(chunkNum uint64) *chunkserver.ChunkServer {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	pickNum := r.Intn(2)
	return m.chunkServers[m.chunksLocation[chunkNum][pickNum]]
}

// 监听rpc信息
// 已测试
func (m *Master) openHeartbeatServer() {
	r := rpc.NewServer()
	r.Register(m)

	addr := fmt.Sprintf(":%v", common.ManagerPort)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	go r.Accept(l)
}

// 更新元数据
// 待完成
// 调用者将更改后的文件传入
// 简化只更改File信息
func (m *Master) UpdateMetaInfo(args *common.UpdateArgs, reply *common.UpdateReply) error {
	for _, file := range args.Files {
		tempFileName := file.FileName
		m.nameSpace.UpdateFile(tempFileName, file)
	}
	return nil
}

// 添加chunkserver
// 待测试
func (m *Master) AddChunkserver(newChunkServer *chunkserver.ChunkServer, chunkserverIp string) {
	m.chunkServers[chunkserverIp] = newChunkServer
}

// 添加目录rpc包装，已测试
func (m *Master) CreateDirectoryRpc(args *common.CreateArgs, reply *common.CreateReply) error {
	//log.Println("rpc--------")
	m.nameSpace.CreateDirectory(args.Test1, args.Test2)
	return nil
}

// 待测试
// 需要追加一个错误信息返回：如果查询不到文件信息，则返回error(已修改)
func (m *Master) FindFileRpc(args *common.GetFileInfoArg, reply *common.GetFileInfoReply) error {
	file, err := m.nameSpace.FindFile(string(args.Path))
	reply.Length = file.FileLength
	reply.Chunks = int64(len(file.Chunks))
	return err
}

func (m *Master) chunkLocations(args *common.GetReplicasArg, reply *common.GetReplicasReply) error {
	locations := m.chunksLocation[uint64(args.Handle)]
	for _, location := range locations {
		reply.Locations = append(reply.Locations, location.Address)
	}
	return nil
}

// 创建文件,根据文件索引创建chunk，输入文件路径和长度
// 待解决，返回值是一个address还是所有address？
func (m *Master) CreateFileRpc(args *common.CreateFileArg, reply *common.CreateFileReply) error {
	//首先创建路径
	f, _ := m.nameSpace.CreateFile(string(args.Path))
	size := args.Length / common.BlockSize
	if args.Length-size*common.BlockSize > 0 {
		size += 1
	}
	//创建相应chunk并更新文件元数据
	for i := 0; i < int(size); i++ {
		newHandle := m.CreateChunk()
		f.Chunks = append(f.Chunks, uint64(newHandle))
		reply.Address = append(reply.Address, m.PickChunkServer(uint64(newHandle)).Address)
		reply.Handle = append(reply.Handle, newHandle)
	}

	return nil

}

// 输入路径，chunk删掉，调用rpcDeleteChunk（传入chunkHandle）
// 待解决，rpcdeldete参数
func (m *Master) DeleteFileRpc(args *common.DeleteFileArg, reply *common.DeleteFileReply) {
	f, _ := m.nameSpace.FindFile(string(args.Path))
	for i := 0; i < len(f.Chunks); i++ {
		for j := 0; j < 3; j++ {
			c, _ := rpc.Dial("tcp", m.chunksLocation[f.Chunks[i]][j])

			defer c.Close()
			var args = &common.DeleteChunkArgs{common.ChunkHandle(f.Chunks[i])}
			var reply = &common.DeleteChunkArgs{}
			c.Call("chunkServer.RPCDeleteChunk", args, reply)
		}

	}
}

// 创建新的chunkHandle,更新元数据，rpc让chunkserver创建chunk
func (m *Master) CreateChunk() common.ChunkHandle {
	newHandle := m.uuid.GenerateUuid() //生成handle
	var randNums []int
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Intn(4)
	var tag [5]int = [5]int{0, 0, 0, 0, 0}
	//挑选3个chunkserver
	for {
		if len(randNums) >= 3 {
			break
		}
		pickNum := r.Intn(4)
		if tag[pickNum] == 0 {
			randNums = append(randNums, pickNum)
			tag[pickNum] = 1
		} else {
			pickNumPlus := (pickNum + 1) % 5
			if tag[pickNumPlus] == 0 {
				randNums = append(randNums, pickNumPlus)
				tag[pickNumPlus] = 1
			}
		}
	}
	//三次rpc
	for i, randNum := range randNums {
		var pickIp string
		switch randNum {
		case 0:
			pickIp = "127.0.0.5:1234"
		case 1:
			pickIp = "127.0.0.6:1234"
		case 2:
			pickIp = "127.0.0.7:1234"
		case 3:
			pickIp = "127.0.0.8:1234"
		case 4:
			pickIp = "127.0.0.9:1234"
		}
		c, _ := rpc.Dial("tcp", pickIp)

		defer c.Close()
		var args = &common.CreateChunkArg{common.ChunkHandle(newHandle)}
		var reply = &common.CreateChunkReply{}
		c.Call("chunkServer.RPCCreateChunk", args, reply)
		//更新元数据
		m.chunksLocation[newHandle][i] = string(pickIp)
	}
	return common.ChunkHandle(newHandle)
}

// 创建chunk
// 待解决args与reply
func (m *Master) createChunkRpc(args *common.CreateChunkRpcArgs, reply *common.CreateChunkRpcReply) error {
	newHandle := m.CreateChunk()
	f, _ := m.nameSpace.FindFile(args.Path)
	reply.Handle = newHandle
	f.Chunks = append(f.Chunks, uint64(newHandle))
	for i := 0; i < len(m.chunksLocation[uint64(newHandle)]); i++ {
		reply.Addresses = append(reply.Addresses, common.ServerAddress(m.chunksLocation[uint64(newHandle)][i]))
	}
	return nil
}

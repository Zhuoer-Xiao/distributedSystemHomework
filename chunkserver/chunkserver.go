package chunkserver

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"path"

	"GFS_Homework/common"
)

// chunkserver的数据结构
type ChunkServer struct {
	Address common.ServerAddress
	master  common.ServerAddress
	rootDir string
	Port    int

	chunk map[common.ChunkHandle]*chunkInfo
}

type chunkInfo struct {
	length common.Offset // 当前chunk偏移量
}

const (
	MetaFileName = "ChunkServer_MetaData"
	FilePerm     = 0755
)

// 创建一个新的ChunkServer
// 已测试
func NewChunkServer(csIP, masterIP common.ServerAddress, rootDir string) *ChunkServer {
	cs := &ChunkServer{
		Address: csIP,
		master:  masterIP,
		rootDir: rootDir,
		chunk:   make(map[common.ChunkHandle]*chunkInfo),
	}
	port := 1234
	cs.Port = port

	//Register
	return cs
}

// 监听rpc信息
func (cs *ChunkServer) HeartBeat() {
	r := rpc.NewServer()
	r.Register(cs)
	addr := fmt.Sprintf("%v", cs.Port)
	l, errx := net.Listen("tcp", addr)
	if errx != nil {
		log.Fatal("ChunkServer Listen Error: ", errx)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}
		go rpc.ServeConn(conn)
	}
}

func (cs *ChunkServer) main() error {
	cs = new(ChunkServer)
	cs.HeartBeat()

	return nil
}

// 元数据存储
// 已测试
func (cs *ChunkServer) StoreMetaData() error {
	filename := path.Join(cs.rootDir, string(cs.Address+MetaFileName))
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, FilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	var metas []common.PersistentChunkInfo
	for handle, ck := range cs.chunk {
		metas = append(metas, common.PersistentChunkInfo{
			Handle: handle, Length: ck.length,
		})
	}
	log.Println("Server stored metadata: ", cs.Address)
	enc := gob.NewEncoder(file)
	err = enc.Encode(metas)

	return nil
}

////////////////////////////
// RPC部分               ///
///////////////////////////

// chunkserver提供的readRPC，读chunk的data并且返回
func (cs *ChunkServer) RPCReadChunk(args *common.ReadChunkArg, reply *common.ReadChunkReply) error {
	handle := args.Handle
	_, chunkinfo := cs.chunk[handle]
	if !chunkinfo {
		return fmt.Errorf("Chunk %v doesn't exist", handle)
	}

	var err error
	reply.Data = make([]byte, args.Length)
	reply.Length, err = cs.readChunk(handle, args.Offset, reply.Data)

	if err != nil {
		return err
	}
	return nil
}

// chunkserver提供的writeRPC
func (cs *ChunkServer) RPCWriteChunk(args *common.WriteChunkArg, reply *common.WriteChunkReply) error {
	handle := args.Handle
	_, chunkinfo := cs.chunk[handle]
	if !chunkinfo {
		return fmt.Errorf("Chunk %v doesn't exist", handle)
	}

	var err error
	err = cs.writeChunk(handle, args.Data, args.Offset)

	if err != nil {
		return err
	}
	return nil
}

// chunkserver提供的appendRPC
func (cs *ChunkServer) RPCAppendChunk(args *common.AppendChunkArg, reply *common.AppendChunkReply) error {
	data := args.Data
	handle := args.Handle

	if len(data) > common.MaxChunkSize {
		return fmt.Errorf("Append data exceeds the max chunk size")
	}

	handle = args.Handle
	ck, chunkinfo := cs.chunk[handle]
	if !chunkinfo {
		return fmt.Errorf("Chunk %v doesn't exist", handle)
	}

	newLen := ck.length + common.Offset(len(data))
	offset := ck.length
	if newLen > common.MaxChunkSize { // 一个chunk装不下
		ck.length = common.MaxChunkSize
		reply.ErrorCode = common.AppendExceedChunkSize
	} else {
		ck.length = newLen
	}
	reply.Offset = offset

	var err error
	err = cs.writeChunk(handle, data, offset)

	if err != nil {
		return err
	}
	return nil
}

// chunkserver提供的创建一个新chunkRPC，给定了chunk handle
// 创建文件时使用，传入文件路径和偏移量
func (cs *ChunkServer) RPCCreateChunk(args *common.CreateChunkArg, reply *common.CreateChunkReply) error {
	fmt.Println("Chunk Server : ", cs.Address, " Create chunk ", args.Handle)

	cs.chunk[args.Handle] = &chunkInfo{
		length: 0,
	}
	filename := path.Join(cs.rootDir, fmt.Sprintf("chunk%v.txt", args.Handle))
	_, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644) // 创建新chunk
	if err != nil {
		return err
	}

	return nil
}

// 文件创建部分函数，待完成
func (cs *ChunkServer) RPCCreateAndWrite(args *common.CreateAndWriteArg, reply *common.CreateAndWriteReply) error {
	data, handle := args.Data, args.Handle
	log.Println("Chunk Server : ", cs.Address, " Create and write chunk ", handle)

	cs.chunk[handle] = &chunkInfo{
		length: common.Offset(len(data)),
	}
	filename := path.Join(cs.rootDir, fmt.Sprintf("chunk%v.txt", handle))
	chunk, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, err = chunk.WriteAt(data, 0)
	if err != nil {
		return err
	}

	return nil
}

// 删除chunk，传入chunkhandle
func (cs *ChunkServer) RPCDeleteChunk(args *common.DeleteChunkArgs, reply *common.DeleteChunkReply) error {
	delete(cs.chunk, args.Handle)

	filename := path.Join(cs.rootDir, fmt.Sprintf("chunk%v.txt", args.Handle))
	err := os.Remove(filename)
	return err
}

/////////////////////////
// 文件操作实现          //
/////////////////////////

// 具体实现从chunk中读取数据
// 返回：读取数据长度，error
func (cs *ChunkServer) readChunk(handle common.ChunkHandle, offset common.Offset, data []byte) (int, error) {
	filename := path.Join(cs.rootDir, fmt.Sprintf("chunk%v.txt", handle))

	f, err := os.Open(filename)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	log.Println("Server %v : read chunk %v at offset: %v, length: %v", cs.Address, handle, offset, len(data))
	return f.ReadAt(data, int64(offset))
}

// 具体实现向chunk中写入数据
// 返回：error
func (cs *ChunkServer) writeChunk(handle common.ChunkHandle, data []byte, offset common.Offset) error {
	ck := cs.chunk[handle]

	newLen := offset + common.Offset(len(data)) // 写入后chunk长度
	if newLen > ck.length {                     // 若超过现chunk长度则修改chunkInfo
		ck.length = newLen
	}

	if newLen > common.MaxChunkSize {
		log.Fatal("New chunk length exceeds the max chunk size")
	}

	log.Println("Server %v : write to chunk %v at offset: %v, length: %v", cs.Address, handle, offset, len(data))
	filename := path.Join(cs.rootDir, fmt.Sprintf("chunk%v.txt", handle))
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, FilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteAt(data, int64(offset))
	if err != nil {
		return err
	}

	return nil
}

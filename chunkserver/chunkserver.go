package chunkserver

import (
	"fmt"
	"os"
	"path"

	"GFS_Homework/common"
)

// chunkserver的数据结构
type ChunkServer struct {
	address common.ServerAddress
	master  common.ServerAddress
	rootDir string

	chunk map[common.ChunkHandle]*chunkInfo
}

type chunkInfo struct {
	length common.Offset
}

const (
	FilePerm = 0755
)

////////////////////////////
// RPC部分               ///
///////////////////////////

// chunkserver提供的readRPC，读chunk的data并且返回
func (cs *ChunkServer) RPCReadChunk(args common.ReadChunkArg, reply *common.ReadChunkReply) error {
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
func (cs *ChunkServer) RPCWriteChunk(args common.WriteChunkArg, reply *common.WriteChunkReply) error {
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

// chunkserver提供的创建一个新chunkRPC，给定了chunk handle
func (cs *ChunkServer) RPCCreateChunk(args common.CreateChunkArg, reply common.CreateChunkReply) error {
	fmt.Println("Chunk Server %v : Create chunk %v", cs.address, args.Handle)

	// 若当前chunk handle号已被占用，则返回错误信息
	if _, ck_info := cs.chunk[args.Handle]; ck_info {
		return fmt.Errorf("Chunk %v already exists", args.Handle)
	}

	cs.chunk[args.Handle] = &chunkInfo{
		length: 0,
	}
	filename := path.Join(cs.rootDir, fmt.Sprintf("chunk%v.chk", args.Handle))
	_, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644) // 创建新chunk
	if err != nil {
		return err
	}

	return nil
}

/////////////////////////
// 文件操作实现          //
/////////////////////////

// 具体实现从chunk中读取数据
// 返回：读取数据长度，error
func (cs *ChunkServer) readChunk(handle common.ChunkHandle, offset common.Offset, data []byte) (int, error) {
	filename := path.Join(cs.rootDir, fmt.Sprintf("chunk%v.chk", handle))

	f, err := os.Open(filename)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	fmt.Println("Server %v : read chunk %v at offset: %v, length: %v", cs.address, handle, offset, len(data))
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
		return fmt.Errorf("New chunk length exceeds the max chunk size")
	}

	fmt.Println("Server %v : write to chunk %v at offset: %v, length: %v", cs.address, handle, offset, len(data))
	filename := path.Join(cs.rootDir, fmt.Sprintf("chunk%v.chk", handle))
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

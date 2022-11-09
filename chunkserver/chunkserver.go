package chunkserver

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
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

type Mutation struct {
	mtype  common.MutationType
	data   []byte
	offset common.Offset
}

type chunkInfo struct {
	length    common.Offset
	mutations []*Mutation // 修改缓冲区

}

// chunkserver提供的readRPC，读chunk的data并且返回
func (cs *ChunkServer) RPCReadChunk(args common.ReadChunkArg, reply *common.ReadChunkReply) error {
	handle := args.Handle
	chunk, chunkinfo := cs.chunk[handle]
	if !chunkinfo {
		return fmt.Errorf("Chunk %v doesn't exist", handle)
	}

	//
	var err error
	reply.Data = make([]byte, args.Length)
	reply.Length, err = cs.readChunk(handle, args.Offset, reply.Data)

	if err != nil {
		return err
	}
	return nil
}

// 具体实现从chunk中读取数据
// 返回：读取数据长度，error
func (cs *ChunkServer) readChunk(handle common.ChunkHandle, offset common.Offset, data []byte) (int, error) {
	filename := path.Join(cs.rootDir, fmt.Sprintf("chunk%v.chk", handle))

	f, err := os.Open(filename)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	log.Infof("Server %v : read chunk %v at offset: %v, length: %v", cs.address, handle, offset, len(data))
	return f.ReadAt(data, int64(offset))
}

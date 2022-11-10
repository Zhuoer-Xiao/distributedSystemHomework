package client

// import (
// 	"fmt"
// 	"math/rand"

// 	log "github.com/Sirupsen/logrus"

// 	"GFS_Homework/common"
// )

// const (
// 	READ   = 0x1
// 	WRITE  = 0x2
// 	RDWR   = 0x3
// 	APPEND = 0x4
// 	CREATE = 0x8
// )

// type FileMode uint32

// // 定义文件结构体
// type File struct {
// 	fd int32 // filedata
// }

// type Client struct {
// 	master common.ServerAddress
// }

// // 文件打开
// func OpenFile(filename string, flag int, perm FileMode) error {

// }

// // 文件读操作
// func (c *Client) Read(path common.Path, offset common.Offset, data []byte) (n int, err error) { // 内存分配data长度的空间
// 	var file common.GetFileInfoReply
// 	// master rpc
// 	errx := common.Call(string(c.master), "Master.RPCGetFileInfo", common.GetFileInfoArg{path}, &file)
// 	if errx != nil {
// 		return -1, errx
// 	}

// 	// 判断读偏移量是否合法
// 	if int64(offset/common.MaxChunkSize) > file.Chunks {
// 		return -1, fmt.Errorf("Read offset exceeds the max file size")
// 	}

// 	pos := 0
// 	for pos < len(data) { // 若跨块，则接着读
// 		index := common.ChunkIndex(offset / common.MaxChunkSize)
// 		chunk_offset := offset % common.MaxChunkSize

// 		var handle common.ChunkHandle
// 		handle, err = c.GetChunkHandle(path, index)
// 		if err != nil {
// 			return
// 		}

// 		var n int
// 		for {
// 			n, err = c.ReadChunk(handle, chunk_offset, data[pos:])
// 			if err == nil {
// 				break
// 			}
// 			log.Warning("Read ", handle, " connection error, please try again: ", err)
// 		}

// 		offset += common.Offset(n)
// 		pos += n // 若跨块，则接着读
// 		if err != nil {
// 			break
// 		}
// 	}

// }

// // 文件写操作
// func (f *File) Write(content []byte) (int, error) {

// }

// // 文件追加写操作
// func (f *File) Append() {

// }

// // 查询文件是否存在
// func (f *File) IsExist(filename string) (bool, error) {

// }

// // 文件删除操作
// func (f *File) Delete() {
// }

// // 文件关闭
// func (f *File) Close() error {
// 	return nil
// }

// ////////////////////////////////////
// //具体功能实现
// ////////////////////////////////////

// // 寻找chunkhandle
// func (c *Client) GetChunkHandle(path common.Path, index common.ChunkIndex) (common.ChunkHandle, error) {
// 	var reply common.GetChunkHandleReply
// 	// master rpc向master传path和index(文件的第几个块)，返回chunkhandle
// 	err := common.Call(string(c.master), "Master.RPCGetChunkHandle", common.GetChunkHandleArg{path, index}, &reply)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return reply.Handle, nil
// }

// // 从给定offset开始阅读文件
// func (c *Client) ReadChunk(handle common.ChunkHandle, offset common.Offset, data []byte) (int, error) {
// 	var readLen int // 欲读数据的长度

// 	// 判断文件是否跨chunk，若跨，则只能读本chunk内的内容
// 	if common.Offset(len(data))+offset < common.MaxChunkSize {
// 		readLen = len(data)
// 	} else {
// 		readLen = int(common.MaxChunkSize - offset)
// 	}

// 	// master rpc返回副本位置信息，给master传chunkhandle，返回一个内含副本位置信息的字符型数组
// 	var l common.GetReplicasReply
// 	err := common.Call(string(c.master), "Master.RPCGetReplicas", common.GetReplicasArg{handle}, &l)
// 	if err != nil {
// 		return 0, common.Error{common.UnknownError, err.Error()}
// 	}
// 	location := l.Locations[rand.Intn(len(l.Locations))] // 随机挑选一个副本读
// 	if len(l.Locations) == 0 {
// 		return 0, common.Error{common.UnknownError, "No replica found"}
// 	}

// 	// chunkserver rpc
// 	var r common.ReadChunkReply
// 	r.Data = data
// 	err = common.Call(string(location), "ChunkServer.RPCReadChunk", common.ReadChunkArg{handle, offset, readLen}, &r)
// 	if err != nil {
// 		return 0, common.Error{common.UnknownError, err.Error()}
// 	}

// 	return r.Length, nil
// }

package client

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"GFS_Homework/common"
)

type Client struct {
	master common.ServerAddress
}

// 返回一个gfs client
func NewClient() *Client {
	return &Client{
		master: "127.0.0.1:1234",
	}
}

// 文件创建
// 客户端传给ChunkServer文件路径和ChunkHandle，chunkserver自行进行分块存储
func (c *Client) Create(path common.Path) error {
	var reply common.CreateFileReply
	// master RPC 创建文件，给master一个路径，返回error
	file, err := os.Open(string(path))
	if err != nil {
		return err
	}
	defer file.Close()

	file_info, _ := file.Stat()
	var k = make([]byte, file_info.Size())
	for {
		h, _ := file.Read(k)
		if h == int(file_info.Size()) {
			break
		}
	}
	file_len := file_info.Size()
	err = common.Call(string(c.master), "Master.RPCCreateFile", common.CreateFileArg{path, file_len}, &reply)
	if err != nil {
		return err
	}

	start := 0
	ck_size := 64 << 10
	for {
		var reply1 common.GetReplicasReply
		handle := reply.Handle[start]
		err = common.Call(string(c.master), "Master.chunkLocations", common.GetReplicasArg{Handle: handle}, &reply1)
		if err != nil {
			return err
		}
		data := k[start*ck_size : int(math.Min(float64((start+1)*ck_size), float64(len(k))))]
		i := 0
		for {
			var cs_IP common.ServerAddress
			cs_IP = reply1.Locations[i]
			var r common.CreateAndWriteReply
			err = common.Call(string(cs_IP), "ChunkServer.RPCCreateAndWrite", common.CreateAndWriteArg{data, handle}, &r)
			if err != nil {
				return err
			}
			i++
			if i == 3 {
				break
			}
		}
		start++
		if start == len(reply.Handle) {
			break
		}
	}

	return nil
}

// 文件删除
func (c *Client) Delete(path common.Path) error {
	var reply common.DeleteFileReply
	// master RPC 删除文件，给master一个路径，返回error
	err := common.Call(string(c.master), "Master.RPCDeleteFile", common.DeleteFileArg{path}, &reply)
	if err != nil {
		return err
	}
	return nil
}

// 判断文件是否存在
func (c *Client) IsExist(path common.Path) (ex bool, err error) {
	var reply common.IsExistReply

	err = common.Call(string(c.master), "Master.FindFileRpc", common.IsExistArg{path}, &reply)
	if err != nil {
		return false, fmt.Errorf("File not exist")
	}

	return true, nil
}

// 文件读操作
func (c *Client) Read(path common.Path, offset common.Offset, data []byte) (n int, err error) { // 内存分配data长度的空间
	var file common.GetFileInfoReply
	// master rpc
	err = common.Call(string(c.master), "Master.FindFileRpc", common.GetFileInfoArg{path}, &file)
	if err != nil {
		return -1, err
	}

	// 判断读偏移量是否合法
	if int64(offset/common.MaxChunkSize) > file.Chunks {
		return -1, fmt.Errorf("Read offset exceeds the max file size")
	}

	pos := 0
	for pos < len(data) { // 若跨块，则接着读
		index := common.ChunkIndex(offset / common.MaxChunkSize)
		chunk_offset := offset % common.MaxChunkSize

		var handle common.ChunkHandle
		handle, err = c.GetChunkHandle(path, index)
		if err != nil {
			return
		}

		var n int
		for {
			n, err = c.ReadChunk(handle, chunk_offset, data[pos:])
			if err == nil {
				break
			}
			fmt.Errorf("Read ", handle, " connection error, please try again: ", err)
		}

		offset += common.Offset(n)
		pos += n // 若跨块，则接着读
		if err != nil {
			break
		}
	}

	if err != nil {
		return pos, nil
	} else {
		return pos, err
	}
}

// 文件写操作
// write和append若返回
func (c *Client) Write(path common.Path, offset common.Offset, data []byte) error {
	var file common.GetFileInfoReply
	err := common.Call(string(c.master), "Master.FindFileRpc", common.GetFileInfoArg{path}, &file)
	if err != nil {
		return err
	}

	if int64(offset/common.MaxChunkSize) > file.Chunks {
		return fmt.Errorf("Write offset exceeds the max file size")
	}

	begin := 0
	for {
		index := common.ChunkIndex(offset / common.MaxChunkSize)
		chunk_offset := offset % common.MaxChunkSize

		handle, err := c.GetChunkHandle(path, index)
		if err != nil {
			return err
		}

		writeMax := int(common.MaxChunkSize - chunk_offset)
		var writeLen int
		if begin+writeMax > len(data) { // 剩余空间够写入
			writeLen = len(data) - begin
		} else { // 否则只能写到chunk结尾
			writeLen = writeMax
		}

		for {
			err = c.WriteChunk(handle, chunk_offset, data[begin:begin+writeLen])
			if err == nil {
				break
			}
			fmt.Errorf("Write ", handle, " connection error, please try again: ", err)
		}
		if err != nil {
			return err
		}

		offset += common.Offset(writeLen) // 在总偏移量上记录已写的数据长度
		begin += writeLen

		if begin == len(data) {
			break
		}
	}

	return nil
}

// 文件追加写操作
// 在record append中， client只是指定数据。GFS在其选定的偏移出将数据至少原子性的加入文件一次，并将偏移返回给client
func (c *Client) Append(path common.Path, data []byte) (offset common.Offset, err error) {
	if len(data) > common.MaxChunkSize {
		return 0, fmt.Errorf("Length of appending data exceeds the max chunk size")
	}

	var file common.GetFileInfoReply
	err = common.Call(string(c.master), "Master.FindFileRpc", common.GetFileInfoArg{path}, &file)
	if err != nil {
		return
	}

	start_chunk := common.ChunkIndex(file.Chunks - 1)
	if start_chunk < 0 {
		start_chunk = 0
	}

	var chunk_offset common.Offset
	for {
		var handle common.ChunkHandle
		handle, err = c.GetChunkHandle(path, start_chunk)
		if err != nil {
			return
		}

		for {
			chunk_offset, err = c.AppendChunk(handle, data)
			if err == nil || err.(common.Error).Code == common.AppendExceedChunkSize { // append内容超出chunk容量，尝试在下一个chunk中继续append
				break
			}
			log.Println("Append ", handle, " connection error, try again ", err)
			time.Sleep(50 * time.Millisecond)
		}
		if err == nil || err.(common.Error).Code != common.AppendExceedChunkSize {
			break
		}

		// 尝试在下一个chunk中append
		start_chunk++
		log.Println("Create a new chunk, Try on the next chunk ", start_chunk)
		var r common.ApplyChunkReply
		common.Call(string(c.master), "Master.RPCApplyNewChunk", common.ApplyChunkArg{path}, &r)
		handle = r.Handle
		chunk_offset, err = c.AppendChunk(handle, data)
		if err == nil {
			break
		}
	}

	offset = common.Offset(start_chunk)*common.MaxChunkSize + chunk_offset
	return offset, nil
}

////////////////////////////////////
//具体功能实现                      //
////////////////////////////////////

// 寻找chunkhandle
func (c *Client) GetChunkHandle(path common.Path, index common.ChunkIndex) (common.ChunkHandle, error) {
	var reply common.GetChunkHandleReply
	// master rpc向master传path和index(文件的第几个块)，返回chunkhandle
	err := common.Call(string(c.master), "Master.OpenFile", common.GetChunkHandleArg{path, index}, &reply)
	if err != nil {
		return 0, err
	}

	return reply.Handle, nil
}

// 从给定offset开始阅读文件
// 返回：读取数据长度，error
func (c *Client) ReadChunk(handle common.ChunkHandle, offset common.Offset, data []byte) (int, error) {
	var readLen int // 欲读数据的长度

	// 判断文件是否跨chunk，若跨，则只能读本chunk内的内容
	if common.Offset(len(data))+offset < common.MaxChunkSize {
		readLen = len(data)
	} else {
		readLen = int(common.MaxChunkSize - offset)
	}

	// master rpc返回副本位置信息，给master传chunkhandle，返回一个内含副本位置信息的字符型数组
	var l common.GetReplicasReply
	err := common.Call(string(c.master), "Master.chunkLocations", common.GetReplicasArg{handle}, &l)
	if err != nil {
		return 0, common.Error{common.UnknownError, err.Error()}
	}
	location := l.Locations[rand.Intn(len(l.Locations))] // 随机挑选一个副本读
	if len(l.Locations) == 0 {
		return 0, common.Error{common.UnknownError, "No replica found"}
	}

	// chunkserver rpc
	var r common.ReadChunkReply
	r.Data = data
	err = common.Call(string(location), "ChunkServer.RPCReadChunk", common.ReadChunkArg{handle, offset, readLen}, &r)
	if err != nil {
		return 0, common.Error{common.UnknownError, err.Error()}
	}

	return r.Length, nil
}

// 从给定offset开始写入文件
func (c *Client) WriteChunk(handle common.ChunkHandle, offset common.Offset, data []byte) error {
	if len(data)+int(offset) > common.MaxChunkSize {
		return fmt.Errorf("Current data lengths + Current offset exceeds the max chunk size")
	}

	var l common.GetReplicasReply
	err := common.Call(string(c.master), "Master.chunkLocations", common.GetReplicasArg{handle}, &l)
	if err != nil {
		return common.Error{common.UnknownError, err.Error()}
	}

	// 不考虑租约，依次写入所有副本
	current := 0
	for {
		location := l.Locations[current]

		if len(l.Locations) == 0 {
			return common.Error{common.UnknownError, "No replica found"}
		}

		// ChunkServer rpc
		var r common.WriteChunkReply
		wcargs := common.WriteChunkArg{handle, offset, data}
		err = common.Call(string(location), "ChunkServer.RPCWriteChunk", wcargs, &r)
		if err != nil {
			return common.Error{common.UnknownError, err.Error()}
		}

		current += 1
		if current == len(l.Locations) {
			break
		}
	}

	return nil
}

// append的data长度不应超过 1/4 chunk大小
func (c *Client) AppendChunk(handle common.ChunkHandle, data []byte) (offset common.Offset, err error) {
	if len(data) > common.MaxChunkSize {
		return 0, common.Error{common.UnknownError, fmt.Sprintf("Lenghth of appending data exceeds the max chunk size, should be within 16 KB")}
	}

	var l common.GetReplicasReply
	err = common.Call(string(c.master), "Master.chunkLocations", common.GetReplicasArg{handle}, &l)
	if err != nil {
		return -1, common.Error{common.UnknownError, err.Error()}
	}

	// 不考虑租约，依次append所有副本
	current := 0
	for {
		location := l.Locations[current]

		if len(l.Locations) == 0 {
			return -1, common.Error{common.UnknownError, "No replica found"}
		}

		// ChunkServer rpc
		var r common.AppendChunkReply
		acargs := common.AppendChunkArg{handle, data}
		err = common.Call(string(location), "ChunkServer.RPCAppendChunk", acargs, &r)
		if err != nil {
			return -1, common.Error{common.UnknownError, err.Error()}
		}
		if r.ErrorCode == common.AppendExceedChunkSize {
			return r.Offset, common.Error{r.ErrorCode, "Append over chunks"}
		}

		current += 1
		if current == len(l.Locations) {
			return r.Offset, nil
		}
	}
}

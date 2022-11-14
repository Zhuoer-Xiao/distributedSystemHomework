package common

type Path string
type Offset int64
type ServerAddress string
type ChunkIndex int
type ChunkHandle uint64

type CreateFileArg struct {
	Path   Path
	Length int64
}

type CreateFileReply struct {
	Address []ServerAddress
	Handle  []ChunkHandle
}

type DeleteFileArg struct {
	Path Path
}

type DeleteFileReply struct{}

type IsExistArg struct {
	Path Path
}

type IsExistReply struct {
	Length int64
	Chunks int64
}

type GetFileInfoArg struct {
	Path Path
}

type GetFileInfoReply struct {
	Length int64
	Chunks int64 // 该文件含有几个chunk
}

type GetChunkHandleArg struct {
	Path  Path
	Index ChunkIndex
}

type GetChunkHandleReply struct {
	Handle ChunkHandle
}

type GetReplicasArg struct {
	Handle ChunkHandle
}

type GetReplicasReply struct {
	Locations []ServerAddress
}

type ReadChunkArg struct {
	Handle ChunkHandle
	Offset Offset
	Length int
}

type ReadChunkReply struct {
	Data      []byte
	Length    int
	ErrorCode ErrorCode
}

type WriteChunkArg struct {
	Handle ChunkHandle
	Offset Offset
	Data   []byte
}

type WriteChunkReply struct {
	Length    int
	ErrorCode ErrorCode
}

type AppendChunkArg struct {
	Handle ChunkHandle
	Data   []byte
}

type AppendChunkReply struct {
	Offset    Offset
	ErrorCode ErrorCode
}

type ApplyChunkArg struct {
	Path Path
}

type ApplyChunkReply struct {
	Handle ChunkHandle
}

type CreateChunkArg struct {
	Handle ChunkHandle
}

type CreateChunkReply struct{}

type CreateAndWriteArg struct {
	Data   []byte
	Handle ChunkHandle
}

type DeleteChunkArgs struct {
	Handle ChunkHandle
}

type DeleteChunkReply struct{}

type CreateAndWriteReply struct{}

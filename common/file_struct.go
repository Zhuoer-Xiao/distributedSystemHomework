package common

type Path string
type Offset int64
type ServerAddress string
type ChunkIndex int
type ChunkHandle uint64

type CreateFileArg struct {
	Path Path
}

type CreateFileReply struct{}

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
	Remain    int64
	ErrorCode ErrorCode
}

type CreateChunkArg struct {
	Handle ChunkHandle
}

type CreateChunkReply struct {
	ErrorCode ErrorCode
}
package common

type Path string
type Offset int64
type ServerAddress string
type ChunkIndex int
type ChunkHandle int64

type ErrorCode int

type Error struct {
	Code ErrorCode
	Err  string
}

type MutationType int

type GetFileInfoArg struct {
	Path Path
}

type GetFileInfoReply struct {
	IsDir  bool
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

type WriteArgs struct {
	Fd  int32
	Off int64
}

const (
	MaxChunkSize = 64 << 10

	UnknownError = -2
)

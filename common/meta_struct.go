package common

// Error
const (
	Success = iota
	UnknownError
	Timeout
	AppendExceedChunkSize
)

const (
	MaxChunkSize = 64 << 10
)

type ErrorCode int
type Checksum int64

type Error struct {
	Code ErrorCode
	Err  string
}

type PersistentChunkInfo struct {
	Handle   ChunkHandle
	Length   Offset
	Checksum Checksum
}

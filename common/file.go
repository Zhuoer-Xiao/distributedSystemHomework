package common

//定义文件
type File struct {
	FileName string   //文件名
	Chunks   []uint64 //所拥有的chunk
}

func NewFile(name string) *File {
	f := new(File)

	f.FileName = name
	f.Chunks = make([]uint64, 0)

	return f
}

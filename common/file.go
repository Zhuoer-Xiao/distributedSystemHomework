package common

import "container/list"

//定义文件
type File struct {
	fileName string     //文件名
	chunks   *list.List //所拥有的chunk
}

func NewFile(name string) *File {
	f := new(File)

	f.fileName = name
	f.chunks = list.New()

	return f
}

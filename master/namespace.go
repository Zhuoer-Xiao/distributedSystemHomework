package master

import (
	"distributedSystemHomework/common"
)

type Directory struct {
	subDir map[string]*Directory   //该路径下的子路径
	files  map[string]*common.File //该路径下的子文件
}

// 新建Directory
func NewDirectory() *Directory {
	return &Directory{make(map[string]*Directory), make(map[string]*common.File)}
}

package master

import (
	"distributedSystemHomework/common"
	"errors"
	"strings"
)

// 目录
type Directory struct {
	subDir map[string]*Directory   //该路径下的子路径
	files  map[string]*common.File //该路径下的子文件
}

// 命名空间
type NameSpace struct {
	rootdir *Directory
}

// 新建Directory
func NewDirectory() *Directory {
	return &Directory{make(map[string]*Directory), make(map[string]*common.File)}
}

// 新建命名空间
func NewNameSpace() *NameSpace {
	return &NameSpace{NewDirectory()}
}

// 递归查找路径
func (d *Directory) recursiveFindDirectory(subpath string) *Directory {
	slice := strings.SplitN(subpath, "/", 2) //按/拆分成两个子字符串
	subdir := d.subDir[slice[0]]
	if subdir != nil {
		return subdir.recursiveFindDirectory(slice[1])
	}
	return nil
}

// 在命名空间查找文件
func (ns *NameSpace) findFile(path string) (*common.File, error) {
	lastSlash := strings.LastIndex(path, "/")
	filename := path
	d := ns.rootdir
	if lastSlash != -1 {
		//slice := strings.Split(path, "/")
		d = ns.rootdir.recursiveFindDirectory(string(path[0:lastSlash]))
		filename = string(path[lastSlash+1:])
		if d == nil {
			return nil, errors.New("No Such File or Directory")
		}
	}

	msg := d.files[filename]
	if msg == nil {
		return nil, errors.New("No Such File")
	}

	return msg, nil
}

func (ns *NameSpace) createFile(path string, flag int, perm uint32) (*common.File, error) {
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash != -1 {
		slice := strings.Split(path, "/")
		d := ns.rootdir.recursiveFindDirectory(string(path[0:lastSlash]))
		if d == nil {
			return nil, errors.New("No Such File of Directory")
		}
		filename := slice[len(slice)-1]
		file := common.NewFile(filename)
		d.files[filename] = file
		return file, nil
	} else {
		file := common.NewFile(path)
		ns.rootdir.files[path] = file
		return file, nil
	}
}

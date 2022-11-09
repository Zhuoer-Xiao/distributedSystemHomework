package master

import (
	"distributedSystemHomework/common"
	"errors"

	//"fmt"
	"strings"
)

// 目录
type Directory struct {
	SubDir map[string]*Directory   //该路径下的子路径
	Files  map[string]*common.File //该路径下的子文件
}

// 命名空间
type NameSpace struct {
	Rootdir *Directory
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
// 已测试
// 输入格式为“example/”，
func (d *Directory) RecursiveFindDirectory(subpath string) *Directory {

	slice := strings.SplitN(subpath, "/", 2) //按/拆分成两个子字符串
	if len(slice) == 1 {
		return d
	}
	//fmt.Println(slice[0])
	subdir := d.SubDir[slice[0]]
	if subdir != nil {
		return subdir.RecursiveFindDirectory(slice[1])
	}
	return nil

}

// 在命名空间查找文件
// 已测试
func (ns *NameSpace) FindFile(path string) (*common.File, error) {
	lastSlash := strings.LastIndex(path, "/")
	filename := path[lastSlash+1:]
	d := ns.Rootdir
	if lastSlash != -1 {
		//slice := strings.Split(path, "/")
		directoryPath := string(path[0 : lastSlash+1])
		d = ns.Rootdir.RecursiveFindDirectory(directoryPath)
		filename = string(path[lastSlash+1:])
		if d == nil {
			return nil, errors.New("No Such File or Directory")
		}
	}

	msg := d.Files[filename]
	if msg == nil {
		return nil, errors.New("No Such File")
	}

	return msg, nil
}

// 输入文件名，如果没有该文件则创建，如果有则返回该文件信息
// 已测试
func (ns *NameSpace) CreateFile(path string) (*common.File, error) {
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash != -1 {
		slice := strings.Split(path, "/")
		d := ns.Rootdir.RecursiveFindDirectory(string(path[0 : lastSlash+1]))
		if d == nil {
			return nil, errors.New("No Such File of Directory")
		}
		filename := slice[len(slice)-1]
		file := common.NewFile(filename)
		d.Files[filename] = file
		return file, nil
	} else {
		file := common.NewFile(path)
		ns.Rootdir.Files[path] = file
		return file, nil
	}
}

// 创建目录
// 已测试
// 输入上级目录和新目录名,如果上级目录为空，不输入"/",否则以"/"结尾
func (ns *NameSpace) CreateDirectory(path string, name string) error {
	newDir := NewDirectory()
	d := ns.Rootdir.RecursiveFindDirectory(path)
	d.SubDir[name] = newDir
	return nil
}

// 删除目录
// 已测试
func (ns *NameSpace) DeleteDirectory(path string, name string) error {
	d := ns.Rootdir.RecursiveFindDirectory(path)
	d.SubDir[name] = nil
	return nil
}

// 删除文件
// 已测试
func (ns *NameSpace) DeleteFile(path string) error {
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash != -1 {
		slice := strings.Split(path, "/")
		d := ns.Rootdir.RecursiveFindDirectory(string(path[0 : lastSlash+1]))
		if d == nil {
			return errors.New("No Such File of Directory")
		}
		filename := slice[len(slice)-1]
		d.Files[filename] = nil
		return nil
	} else {
		ns.Rootdir.Files[path] = nil
		return nil
	}
}

// 更新文件，
// 已测试
func (ns *NameSpace) UpdateFile(path string, file *common.File) error {
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash != -1 {
		slice := strings.Split(path, "/")
		d := ns.Rootdir.RecursiveFindDirectory(string(path[0 : lastSlash+1]))
		if d == nil {
			return errors.New("No Such File of Directory")
		}
		filename := slice[len(slice)-1]
		d.Files[filename] = file
		return nil
	} else {
		ns.Rootdir.Files[path] = file
		return nil
	}
}

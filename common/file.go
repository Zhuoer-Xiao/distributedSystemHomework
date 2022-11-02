package common

import "container/list"

//定义文件
type File struct{
	fileName string//文件名
	chunks *list.List//所拥有的chunk
}
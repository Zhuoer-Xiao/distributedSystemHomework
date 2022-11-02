package common
import(

)
type FileToChunk struct{

}

//master元数据
// 文件命名空间
// 文件到数据块的映射信息
// 数据块的位置信息
// 访问控制信息
// 数据块版本号
type MetaData struct{
	fileNameSpace string
	fileToChunk map[string]*[]uint64
	chunkInChunkServer map[uint64]uint64
	chunkVersion uint64
}
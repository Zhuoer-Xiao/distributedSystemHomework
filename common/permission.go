package common
import(

)
//访问控制
//通过位运算实现
const(
	READ =0x1
	WRITE =0x2
	APPEND = 0x4
	CREATE =0x8
	EXIST =0x16
	DELETE =0x32
)

func CheckRead(perm uint32) bool{
	return ((perm & 0x1)!=0)
}

func CheckWrite(perm uint32) bool{
	return ((perm & 0x2)!=0)
}

func CheckAppend(perm uint32) bool{
	return ((perm & 0x4)!=0)
}

func CheckCreate(perm uint32) bool{
	return ((perm & 0x8)!=0)
}

func CheckExist(perm uint32) bool{
	return ((perm & 0x16)!=0)
}

func CheckDelete(perm uint32) bool{
	return ((perm & 0x32)!=0)
}


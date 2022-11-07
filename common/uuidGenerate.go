package common
//生成唯一标识符
//待测试

type UuidGenerate struct{
	UuidNow uint64
}

func(u *UuidGenerate)NewUuidGnerate(){
	u.UuidNow=0
}

func (u *UuidGenerate)GenerateUuid()uint64{
	u.UuidNow++
	return u.UuidNow
}
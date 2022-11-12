package common
//生成唯一标识符
//待测试

import(
	"sync"
)

var Mu sync.Mutex
type UuidGenerate struct{
	UuidNow uint64
	mu sync.Mutex
}

func NewUuidGnerate()*UuidGenerate{
	Uu:=&UuidGenerate{}
	Uu.UuidNow=0
	return Uu
}

func (u *UuidGenerate)GenerateUuid()uint64{
	u.mu.Lock()
	u.UuidNow++
	defer u.mu.Unlock()
	return u.UuidNow
}
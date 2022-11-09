package master

import (
	"distributedSystemHomework/common"
	"fmt"
)
func TestMaster(){
	m:=NewMaster()
	m.nameSpace.CreateDirectory("","node1")
	m.nameSpace.CreateDirectory("node1/","node2")
	//m.nameSpace.DeleteDirectory("node1/","node2")
	m.nameSpace.CreateFile("node1/node2/test1.txt")
	file:=common.NewFile("test2.txt")
	m.nameSpace.UpdateFile("node1/node2/test1.txt",file)
	fmt.Println("---------------------")
}
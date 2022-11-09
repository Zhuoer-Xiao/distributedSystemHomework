package main

import (
	"distributedSystemHomework/common"
	"distributedSystemHomework/master"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	// "distributedSystemHomework/chunkserver"
)

func main2() {
	myDir1 := master.NewDirectory()
	myDir2 := master.NewDirectory()
	myDir3 := master.NewDirectory()
	myDir4 := master.NewDirectory()
	myDir5 := master.NewDirectory()
	myName := master.NewNameSpace()
	myName.Rootdir.SubDir["node1"] = myDir1
	myDir1.SubDir["node2"] = myDir2
	myDir1.SubDir["node3"] = myDir3
	myDir2.SubDir["node4"] = myDir4
	myDir2.SubDir["node5"] = myDir5
	myFile := common.NewFile("test1.txt")
	myDir3.Files["test1.txt"] = myFile
	path := "node1/node3/test1.txt"

	lastSlash := strings.LastIndex(path, "/")
	fmt.Println(string(path[0:lastSlash]))
	fmt.Println(string(path[lastSlash+1:]))
	fmt.Println("---------------------------")

	res, _ := myName.CreateFile("node1/node3/file.txt")

	file1, _ := os.OpenFile("./file1.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	outPut1, _ := json.Marshal(myName)
	file1.Write(outPut1)
	file1.Close()


	fmt.Println(res.FileName)
}

func main(){
	master.TestMaster()
	fmt.Println("-------------------------")
}

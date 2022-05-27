package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

type list struct {
	startCHS string
	startCBA string
	size     float64
	Type     string
}

func main() {
	filename := ""
	if len(os.Args) > 1 {
		filename = strings.Join(os.Args[1:], "")
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("文件打开失败", err.Error())
		return
	}
	defer file.Close()
	var chunk []byte
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil && err != io.EOF {
		fmt.Println("read buf fail", err)
		return
	}
	chunk = buf[446:]
	fmt.Println(chunk)
	fmt.Println(len(chunk))
	convInt := binary.LittleEndian.Uint32(chunk[12:16])
	fmt.Println(convInt)

}

func getResult(chunk []byte) list {
	var c, h, s int
	result := list{}
	c = int(chunk[3])
	h = int(chunk[1])
	s = int(chunk[2])
	result.startCBA = fmt.Sprintf("C:%d,H:%d,:%d", c, h, s)

	return result
}

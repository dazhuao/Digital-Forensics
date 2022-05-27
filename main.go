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
	startLBA uint32
	size     float32
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
	chunk = buf[446:462]
	result := getResult(chunk)
	fmt.Println(result)

}

func getResult(chunk []byte) list {
	var c, h, s int
	result := list{}
	c = int(chunk[3])
	h = int(chunk[1])
	s = int(chunk[2])
	result.startCHS = fmt.Sprintf("C:%d,H:%d,S:%d", c, h, s)
	result.startLBA = binary.LittleEndian.Uint32(chunk[8:12])
	numberOfs := int(binary.LittleEndian.Uint32(chunk[12:]))
	result.size = float32(numberOfs) * 512 / 1024 / 1024
	switch chunk[4] {
	case 1:
		result.Type = "FAT12"
	case 4, 6:
		result.Type = "FAT16"
	case 5:
		result.Type = "Extend"
	case 7:
		result.Type = "NTFS"
	case 11, 12:
		result.Type = "FAT32"
	default:
		result.Type = ""
	}
	return result
}

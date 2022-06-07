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
	var resultList []list
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

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read buf fail", err)
			return
		}
		if n == 0 {
			break
		}
		chunk = append(chunk, buf[:n]...)
	}
	resultList = getResultList(chunk, 0)
	fmt.Println("Seq.# Starting CHS Starting LBA Size (MB) Type")
	for n, list := range resultList {
		fmt.Println(n+1, list.startCHS, list.startLBA, list.size, list.Type)
	}

}

// input 16byte, out put the list
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

func isnull(chunk []byte) bool {
	for _, s := range chunk {
		if s != 0 {
			return false
		}
	}
	return true
}

func getResultList(chunk []byte, org int) []list {
	var resultList []list
	for n := 0; n <= 3; n++ {
		MBR := chunk[org+446+n*16 : org+462+n*16]
		if isnull(MBR) {
			continue
		}
		result := getResult(MBR)
		resultList = append(resultList, result)
		if result.Type == "Extend" {
			resultList2 := getResultList(chunk, int(result.startLBA*512))
			for m := range resultList2 {
				resultList2[m].startLBA = resultList2[m].startLBA + result.startLBA
			}
			resultList = append(resultList, resultList2...)
		}
	}
	return resultList
}

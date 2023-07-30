package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main01() {
	// 打开一个100字节的文件
	file, err := os.Open("./main.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 设置每次读取的字节数
	chunkSize := 78

	// 设置缓冲区
	buffer := make([]byte, chunkSize)

	// 循环10次，每次从不同的偏移量开始读取
	for i := 0; i < 10; i++ {
		// 计算偏移量
		offset := int64(i * chunkSize)

		// 从偏移量开始读取数据到缓冲区
		bytesRead, err := file.ReadAt(buffer, offset)
		fmt.Printf("%s", buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("已经读完文件", bytesRead)
				break
			} else {
				log.Fatal(err)
			}
		}

		// 打印读取的数据和字节数
		fmt.Printf("从偏移量 %d 开始读取了 %d 个字节: %s\n", offset, bytesRead, buffer[:bytesRead])
	}
}
func main() {

	t, _ := os.Open(".\\fs\\user\\1\\test.exe")
	h, _ := os.Create(".\\fs\\user\\1\\test2.exe")
	buf := make([]byte, 1000)
	var offset int64
	for {
		n1, err1 := t.ReadAt(buf, offset)
		offset += int64(n1)
		_, _ = h.Write(buf)
		if err1 != nil && err1 != io.EOF {
			fmt.Println(err1)
		}
		fmt.Println(n1)
		if err1 == io.EOF {
			break
		}

	}

}

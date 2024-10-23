package main

import (
	"ffmpeg_hls_go/pkg/sm3"
	"fmt"
)

func main() {
	message := []byte("Hello, SM3!")

	// 创建一个 SM3 哈希对象
	hash := sm3.New()

	// 写入消息数据
	hash.Write(message)

	// 计算哈希值
	digest := hash.Sum(nil)

	// 打印哈希结果
	fmt.Printf("SM3 Hash: %x\n", digest)
}

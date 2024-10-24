package main

import (
	"ffmpeg_hls_go/internal/logger"
	"ffmpeg_hls_go/internal/token"
	"ffmpeg_hls_go/pkg/sm3"
	"time"
)

func main() {
	log := logger.GetLoggerInstance()
	log.Info("xxxx")
	message := []byte("Hello, SM3!")

	// 创建一个 SM3 哈希对象
	hash := sm3.New()

	// 写入消息数据
	hash.Write(message)

	// 计算哈希值
	digest := hash.Sum(nil)

	// 打印哈希结果
	log.Printf("SM3 Hash: %x", digest)

	token.Init()
	ch := make(chan bool)
	for {
		log.Println(token.GetResponse().AccessToken)
		time.Sleep(1 * time.Second)
	}
	<-ch
	log.Println("done")
}

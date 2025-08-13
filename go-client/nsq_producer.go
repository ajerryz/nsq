package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

func main() {
	// 创建配置
	config := nsq.NewConfig()
	config.ReadTimeout = 40 * time.Second
	config.WriteTimeout = 40 * time.Second

	// 创建生产者
	producer, err := nsq.NewProducer("localhost:4150", config)
	if err != nil {
		panic(fmt.Errorf("创建 Producer 失败，err:%w", err))
	}
	defer producer.Stop()

	// 测试发送消息
	if err := producer.Publish("demo-topic", []byte("hello world")); err != nil {
		panic(fmt.Errorf("publish 消息失败,err:%w", err))
	}
	fmt.Println("Publish Message Success")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

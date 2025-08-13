package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

// MyHandler 自定义消息处理函数
type MyHandler struct{}

func (h *MyHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	log.Println("Received message: ", string(m.Body))
	return nil
}

func main() {
	// 创建配置
	cfg := nsq.NewConfig()
	cfg.MsgTimeout = 10 * time.Second
	// 创建消费者
	consumer, err := nsq.NewConsumer("demo-topic", "demo-channel", cfg)
	if err != nil {
		panic(fmt.Errorf("创建消费者失败,err:%w", err))
	}
	// 添加消息处理器
	consumer.AddHandler(&MyHandler{})

	// 连接到NSQD 节点 或者 NSQLookupd 服务
	err = consumer.ConnectToNSQD("127.0.0.1:4150") // 方式1 连接到NSQD(适合NSQD单点)
	//err = consumer.ConnectToNSQLookupd("127.0.0.1:4161") // 方式2 连接到NSQLookupd 发现节点(适合集群)
	if err != nil {
		panic(fmt.Errorf("连接失败,err:%w", err))
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	consumer.Stop()

}

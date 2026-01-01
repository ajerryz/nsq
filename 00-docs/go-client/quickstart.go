package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nsqio/go-nsq"
)

const (
	// 消费者通过lookupd 寻找nsqd消费
	nsqlookupdHttpAddr = "localhost:4161"
	// nsqdTcpAddr 生产者直连nsqd生产
	nsqdTcpAddr = "localhost:4150"
	topic       = "goquickstart"
	channel1    = "ch1"
	channel2    = "ch2"
)

func main() {
	go startConsumer()
	startProducer()

	//time.Sleep(100 * time.Second)
}

type SimpleHandler struct {
	log *log.Logger
}

func NewSimpleHandler() *SimpleHandler {
	return &SimpleHandler{
		log: log.New(os.Stdout, "[SimpleHandler]", log.LstdFlags|log.Lmsgprefix),
	}

}

func (h *SimpleHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		h.log.Println("Empty message")
		return nil
	}
	h.log.Printf("%s\n", string(m.Body))
	return nil
}

func startConsumer() {
	addr := nsqlookupdHttpAddr
	config := nsq.NewConfig()
	config.MsgTimeout = 10 * time.Second
	consumer, err := nsq.NewConsumer(topic, channel1, config)
	if err != nil {
		panic(fmt.Errorf("create consumer error: %w", err))
	}
	consumer.SetLoggerLevel(nsq.LogLevelInfo)
	consumer.AddHandler(NewSimpleHandler())
	if err = consumer.ConnectToNSQLookupd(addr); err != nil {
		panic(fmt.Errorf("connect to nsqlookupd error: %w", err))
	}

}

func startProducer() {
	addr := nsqdTcpAddr
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		panic(fmt.Errorf("create producer failed, err: %s", err))
	}
	producer.SetLoggerLevel(nsq.LogLevelInfo)
	if err = producer.Ping(); err != nil {
		panic(fmt.Errorf("ping producer failed, err: %s", err))
	}
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("输入要发送的消息(topic:%s):", topic)
		line, _ := r.ReadString('\n')
		fmt.Printf("输入的是:%v\n", line)
		err := producer.Publish(topic, []byte(line))
		if err != nil {
			fmt.Printf("发送失败:%v\n", err)
		} else {
			fmt.Printf("发送成功\n")
		}
	}
}

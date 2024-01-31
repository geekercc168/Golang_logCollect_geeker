package main

import (
	"github.com/Shopify/sarama"
	"logTransfer/global"
	"strings"
	"sync"
)

type KafkaClient struct {
	Client sarama.Consumer
	Addr   string
	Topic  string
	wg     sync.WaitGroup
}

var (
	kafkaClient *KafkaClient
)

func InitKafka(addr string, topic string) (err error) {

	kafkaClient = &KafkaClient{}
	consumer, err := sarama.NewConsumer(strings.Split(addr, "|"), nil)
	if err != nil {
		global.Log.Error("启动Kafka消费者错误: %s", err)
		return nil
	}
	kafkaClient.Client = consumer
	kafkaClient.Addr = addr
	kafkaClient.Topic = topic
	return
}

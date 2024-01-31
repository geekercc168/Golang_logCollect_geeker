package main

import (
	"github.com/Shopify/sarama"
	"logTransfer/es"
	"logTransfer/global"
	"sync"
)

func run() (err error) {
	// Kafka消费数据
	partitionList, err := kafkaClient.Client.Partitions(kafkaClient.Topic)
	if err != nil {
		global.Log.Error("Failed to get the list of partitions: ", err)
		return
	}
	// 用于等待所有协程完成的 WaitGroup
	var wg sync.WaitGroup
	wg.Add(len(partitionList))

	for partition := range partitionList {
		//sarama.OffsetNewest 最新消息
		//sarama.OffsetOldest 已经消费的旧消息
		pc, errRet := kafkaClient.Client.ConsumePartition(kafkaClient.Topic, int32(partition), sarama.OffsetOldest)
		if errRet != nil {
			err = errRet
			global.Log.Error("Failed to start consumer for partition %d: %s\n", partition, err)
			return
		}

		go func(pc sarama.PartitionConsumer) {

			defer pc.AsyncClose() //关闭消费者

			for msg := range pc.Messages() {
				global.Log.Debug("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				err = es.SendToES(kafkaClient.Topic, msg.Value)
				if err != nil {
					global.Log.Warn("send to es failed, err:%v", err)
				}
			}
			// 每个消费者协程完成时通知 WaitGroup
			wg.Done()
		}(pc)
	}

	// 主协程等待所有消费者协程完成
	wg.Wait()
	return
}

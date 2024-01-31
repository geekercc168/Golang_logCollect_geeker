package main

import (
	"logAgent/global"
	"logAgent/kafka"
	"logAgent/tailf"
	"time"
)

func serverRun() (err error) {
	for {
		msg := tailf.GetOneLine()
		err = sendToKafka(msg)
		if err != nil {
			global.Log.Error("Send to Kafka failed, err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}

}

func sendToKafka(msg *tailf.TextMsg) (err error) {
	//fmt.Printf("读取 msg:%s, topic:%s\n", msg.Msg, msg.Topic)
	_ = kafka.SendToKafka(msg.Msg, msg.Topic)
	return
}

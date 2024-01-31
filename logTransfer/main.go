package main

import (
	"github.com/astaxie/beego/logs"
	"logTransfer/es"
	"logTransfer/global"
)

var err error

func main() {

	// 初始化配置
	err = InitConfig("ini", "C:\\Users\\geeker_yuyu\\Desktop\\coder\\workspace\\Golang_logCollect-master\\logTransfer\\config\\logTransfer.configs")
	if err != nil {
		panic(err)
		return
	}
	logs.Info("初始化配置成功")
	//初始化日志模块
	err = initLogger(logConfig.LogPath, logConfig.LogLevel)
	if err != nil {
		panic(err)
		return
	}
	global.Log.Debug("初始化日志模块成功")

	// 初始化Kafka
	err = InitKafka(logConfig.KafkaAddr, logConfig.KafkaTopic)
	if err != nil {
		global.Log.Error("初始化Kafka失败, err:", err)
		return
	}
	global.Log.Debug("初始化Kafka成功")

	// 初始化Es
	err = es.InitEs(logConfig.EsAddr)
	if err != nil {
		global.Log.Error("初始化Elasticsearch失败, err:", err)
		return
	}
	global.Log.Debug("初始化Es成功")

	// 运行
	err = run()
	if err != nil {
		global.Log.Error("运行错误, err:", err)
		return
	}

	global.Log.Warn("logTransfer 退出")
}

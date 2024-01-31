package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logAgent/global"
	"logAgent/kafka"
	"logAgent/tailf"
)

// 定义初始变量
var (
	err error
)

// 初始化
func init() {

	fmt.Println("配置文件开始初始化....")
	// 读取初始化配置文件
	filename := "./configs/logAgent.conf"
	err = loadInitConf("ini", filename)
	if err != nil {
		fmt.Printf("配置文件初始化失败:%v\n", err)
		return
	}
	// 输出成功信息
	fmt.Println("配置文件初始化成功！")

	fmt.Println("日志设置初始化....")
	// 初始化日志信息
	err = initLogger()
	if err != nil {
		fmt.Printf("导入日志文件错误:%v\n", err)
		global.Log.Error(fmt.Sprintf("日志初始化设置失败:%s\r\n", err))
		return
	}
	fmt.Println("日志设置初始化成功！")

	// 初识化etcd
	fmt.Println("etcd开始初始化....")
	collectConf, err := initEtcd(global.LogConfig.EtcdAddr, global.LogConfig.EtcdKey)

	if err != nil {
		fmt.Printf("导入日志文件错误:%v\n", err)
		global.Log.Error(fmt.Sprintf("初始化etcd失败:%s\r\n", err))
		return
	}
	fmt.Println("etcd初始化成功！")

	// 初始化tailf
	err = tailf.InitTail(collectConf, global.LogConfig.ChanSize)
	if err != nil {
		logs.Error("初始化tailf失败:", err)
		return
	}
	global.Log.Debug("初始化tailf成功!")

	// 初始化Kafka
	err = kafka.InitKafka(global.LogConfig.KafkaAddr)
	if err != nil {
		logs.Error("初识化Kafka producer失败:", err)
		return
	}
	global.Log.Debug("初始化Kafka成功!")
}
func main() {
	// 运行
	err = serverRun()
	if err != nil {
		logs.Error("serverRun failed:", err)
	}
	logs.Info("程序退出")
}

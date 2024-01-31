package main

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"logAgent/global"
	"logAgent/tailf"
)

// 日志配置
type Config struct {
	logLevel string
	logPath  string

	chanSize    int
	KafkaAddr   string
	collectConf []tailf.CollectConf

	etcdAddr string
	etcdKey  string
}

// 导入初始化配置
func loadInitConf(confType, filename string) (err error) {
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		fmt.Printf("初始化配置文件出错:%v\n", err)
		return
	}
	// 导入配置信息
	global.LogConfig = &global.Config{}
	// 日志级别
	global.LogConfig.LogLevel = conf.String("logs::log_level")
	if len(global.LogConfig.LogLevel) == 0 {
		global.LogConfig.LogLevel = "debug"
	}
	// 日志输出路径
	global.LogConfig.LogPath = conf.String("logs::log_path")
	if len(global.LogConfig.LogPath) == 0 {
		global.LogConfig.LogPath = "./logs/my.log"
	}

	// 管道大小
	global.LogConfig.ChanSize, err = conf.Int("logs::chan_size")
	if err != nil {
		global.LogConfig.ChanSize = 100
	}

	// Kafka
	global.LogConfig.KafkaAddr = conf.String("kafka::server_addr")
	if len(global.LogConfig.KafkaAddr) == 0 {
		err = fmt.Errorf("初识化Kafka失败")
		return
	}

	// etcd
	global.LogConfig.EtcdAddr = conf.String("etcd::addr")
	if len(global.LogConfig.EtcdAddr) == 0 {
		err = fmt.Errorf("初识化etcd addr失败")
		return
	}

	global.LogConfig.EtcdKey = conf.String("etcd::configKey")
	if len(global.LogConfig.EtcdKey) == 0 {
		err = fmt.Errorf("初识化etcd configKey失败")
		return
	}

	return
}

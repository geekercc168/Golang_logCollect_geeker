package global

import (
	"github.com/astaxie/beego/logs"
	"logAgent/tailf"
)

type Config struct {
	LogLevel string
	LogPath  string

	ChanSize    int
	KafkaAddr   string
	CollectConf []tailf.CollectConf

	EtcdAddr string
	EtcdKey  string
}

var (
	LogConfig *Config
	Log       *logs.BeeLogger
	err       error
)

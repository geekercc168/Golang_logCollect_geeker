package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"logAgent/global"
)

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

func initLogger() (err error) {

	config := make(map[string]interface{})
	config["filename"] = global.LogConfig.LogPath
	config["level"] = convertLogLevel(global.LogConfig.LogLevel)
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("初始化日志, 序列化失败:", err)
		return
	}
	global.Log = logs.NewLogger(10000) //设置日志大小 如果超过 会把日志转移到分割文件中
	global.Log.EnableFuncCallDepth(true)
	global.Log.SetLogFuncCallDepth(3)
	global.Log.SetLogger(logs.AdapterConsole)
	_ = global.Log.SetLogger(logs.AdapterFile, string(configStr))

	return
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"logTransfer/global"
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

func initLogger(logPath string, logLevel string) (err error) {

	config := make(map[string]interface{})
	config["filename"] = logPath
	config["level"] = convertLogLevel(logLevel)
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("初始化日志, 序列化失败:", err)
		return
	}
	global.Log = logs.NewLogger(10000)
	global.Log.EnableFuncCallDepth(true)
	global.Log.SetLogFuncCallDepth(3)
	//log.SetLogger(logs.AdapterConsole)
	_ = global.Log.SetLogger(logs.AdapterFile, string(configStr))

	return
}

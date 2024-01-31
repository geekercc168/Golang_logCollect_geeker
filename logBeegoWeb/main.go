package main

import (
	"fmt"
	"logBeegoWeb/model"
	_ "logBeegoWeb/router"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/coreos/etcd/clientv3"
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func initLog() (err error) {
	logs.SetLevel(7)
	logs.SetPrefix("logBeegoWeb")
	logs.SetLogFuncCall(true)
	//logs.SetLogFuncCallDepth(3)
	err = logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/log.txt"}`)
	if err != nil {
		logs.Warn("initDb failed, err:%v", err)
		return err
	}
	return nil
}
func initDb() (err error) {
	database, err := sqlx.Open("mysql", "root:123456@tcp(127.0.0.1:3309)/logCollect")
	if err != nil {
		logs.Warn("open mysql failed,", err)
		return
	}

	model.InitDb(database)
	return
}

func initEtcd() (err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:12379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	model.InitEtcd(cli)
	return
}

func main() {
	err := initLog()
	if err != nil {
		logs.Warn("init etcd failed, err:%v", err)
		return
	}

	err = initDb()
	if err != nil {
		logs.Warn("initDb failed, err:%v", err)
		return
	}

	err = initEtcd()
	if err != nil {
		logs.Warn("init etcd failed, err:%v", err)
		return
	}

	beego.Run()
}

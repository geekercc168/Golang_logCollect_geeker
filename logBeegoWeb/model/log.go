package model

import (
	"context"
	"encoding/json"
	"logBeegoWeb/services/tailf"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/clientv3"
)

var (
	etcdClient *clientv3.Client
)

type LogInfo struct {
	AppId      int    `db:"app_id"`
	AppName    string `db:"app_name"`
	LogId      int    `db:"log_id"`
	Status     int    `db:"status"`
	CreateTime string `db:"create_time"`
	LogPath    string `db:"log_path"`
	Topic      string `db:"topic"`
}

func InitEtcd(client *clientv3.Client) {
	etcdClient = client
}

func GetAllLogInfo() (loglist []LogInfo, err error) {
	err = Db.Select(&loglist,
		"select a.app_id, b.app_name, a.create_time, a.topic, a.log_id, a.status, a.log_path from tbl_log_info a, tbl_app_info b where a.app_id=b.app_id")
	if err != nil {
		logs.Warn("获取app列表信息失败, err:%v", err)
		return
	}
	return
}

func CreateLog(info *LogInfo) (err error) {

	conn, err := Db.Begin()
	if err != nil {
		logs.Warn("CreateApp failed, Db.Begin error:%v", err)
		return
	}

	defer func() {
		if err != nil {
			_ = conn.Rollback()
			return
		}

		_ = conn.Commit()
	}()

	var appId []int
	err = Db.Select(&appId, "select app_id from tbl_app_info where app_name=?", info.AppName)
	if err != nil || len(appId) == 0 {
		logs.Warn("select app_id failed, Db.Exec error:%v", err)
		return
	}

	info.AppId = appId[0]
	r, err := conn.Exec("insert into tbl_log_info(app_id, log_path, topic)values(?, ?, ?)",
		info.AppId, info.LogPath, info.Topic)

	if err != nil {
		logs.Warn("CreateApp failed, Db.Exec error:%v", err)
		return
	}

	_, err = r.LastInsertId()
	if err != nil {
		logs.Warn("CreateApp failed, Db.LastInsertId error:%v", err)
		return
	}

	return
}

func SetLogConfToEtcd(etcdKey string, info *LogInfo) (err error) {

	var logConfArr []tailf.CollectConf
	logConfArr = append(
		logConfArr,
		tailf.CollectConf{
			LogPath: info.LogPath,
			Topic:   info.Topic,
		},
	)

	data, err := json.Marshal(logConfArr)
	if err != nil {
		logs.Warn("marshal failed, err:%v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//cli.Delete(ctx, EtcdKey)
	//return
	_, err = etcdClient.Put(ctx, etcdKey, string(data))
	cancel()
	if err != nil {
		logs.Warn("Put failed, err:%v", err)
		return
	}

	logs.Debug("put etcd succ, data:%v", string(data))
	return
}

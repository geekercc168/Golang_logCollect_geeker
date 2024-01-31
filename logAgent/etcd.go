package main

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/client/v3"
	"logAgent/global"
	"logAgent/tailf"
	"strings"
	"time"
)

type EtcdClient struct {
	client *clientv3.Client
	keys   []string
}

var (
	etcdClient *EtcdClient
)

// etcd 服务注册与发现 demo 参考资料:https://juejin.cn/post/7101947466722836487
func initEtcd(addr string, etcdKey string) (collectConf []tailf.CollectConf, err error) {

	//分割addr为切片
	addr_str := strings.Split(addr, "|")

	// 初始化连接etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: addr_str,
		//Endpoints: []string{"127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"},
		//Endpoints: []string{addr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("连接etcd失败:", err)
		return
	}

	etcdClient = &EtcdClient{
		client: cli,
	}

	// 如果Key不是以"/"结尾, 则自动加上"/"
	if strings.HasSuffix(etcdKey, "/") == false {
		etcdKey = etcdKey + "/"
	}

	//获取固定键的值
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, etcdKey)
	if err != nil {
		global.Log.Error("etcd get请求失败:", err)
	}
	defer cancel()

	//
	for _, v := range resp.Kvs {
		if string(v.Key) == etcdKey {
			global.Log.Debug("resp from etcd v.Value:%s", v.Value)
			// 反序列化为结构体
			err = json.Unmarshal(v.Value, &collectConf)
			if err != nil {
				global.Log.Error("反序列化失败:", err)
				continue
			}
			global.Log.Debug("日志设置为%v", collectConf)
		}
	}
	initEtcdWatcher(global.LogConfig.EtcdAddr)
	global.Log.Debug("连接etcd成功")
	return
}

// 初始化多个watch监控etcd中配置节点
func initEtcdWatcher(addr string) {
	for _, key := range etcdClient.keys {
		go watchKey(addr, key)
	}
}

func watchKey(addr string, key string) {

	//分割addr为切片
	addr_str := strings.Split(addr, "|")

	// 初始化连接etcd
	cli, err := clientv3.New(clientv3.Config{
		//Endpoints: []string{"127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"},
		Endpoints:   addr_str,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		global.Log.Error("连接etcd失败:", err)
		return
	}

	global.Log.Debug("开始监控key:", key)

	// Watch操作
	var collectConf []tailf.CollectConf
	var getConfSucc = true
	wch := cli.Watch(context.Background(), key)
	for resp := range wch {
		for _, ev := range resp.Events {
			// DELETE处理
			//mvccpb.DELETE
			if ev.Type.String() == "1" {
				logs.Warn("删除Key[%s]配置", key)
				continue
			}
			// PUT处理
			//mvccpb.PUT
			if ev.Type.String() == "0" && string(ev.Kv.Key) == key {
				err = json.Unmarshal(ev.Kv.Value, &collectConf)
				if err != nil {
					logs.Error("反序列化key[%s]失败:", err)
					getConfSucc = false
					continue
				}
			}
			global.Log.Debug("get config from etcd ,Type: %v, Key:%v, Value:%v\n", ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
		}

		if getConfSucc {
			global.Log.Debug("get config from etcd success, %v", collectConf)
			_ = tailf.UpdateConfig(collectConf)
		}
	}
}

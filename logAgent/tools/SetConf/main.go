package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"logAgent/tailf"
	"time"
)

// 定义etcd的前缀key
const (
	EtcdKey = "/backend/logagent/config/127.0.0.1"
)

func SetLogConfToEtcd() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")
	defer cli.Close()

	var logConfArr []tailf.CollectConf
	logConfArr = append(
		logConfArr,
		tailf.CollectConf{
			LogPath: "../../logs/my.log",
			Topic:   "nginx_log",
		},
	)

	// Json打包
	data, err := json.Marshal(logConfArr)
	if err != nil {
		fmt.Println("json failed, ", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, EtcdKey, string(data))
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}

func main() {
	SetLogConfToEtcd()
}

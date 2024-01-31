package es

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/olivere/elastic/v7"
	"logTransfer/global"
	"net/http"
)

type LogMessage struct {
	App     string
	Topic   string
	Message string
}

var (
	esClient *elastic.Client
)

const mapping = `
{
  "mappings": {
      "properties":{ 
		  "message":{
			  "type":"text",
			  "fields": {
				  "keyword": {
					"type": "keyword", 
					"ignore_above": 256
				  }
			  }
		  }
	  }
  }
}`

func InitEs(addr string) (err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client_https := &http.Client{Transport: tr} // 自定义transport

	client, err := elastic.NewClient(
		elastic.SetHttpClient(client_https),
		elastic.SetSniff(false),
		elastic.SetURL(addr),
		elastic.SetBasicAuth("elastic", "123456"), //设置账号、密码
	)
	if err != nil {
		fmt.Println("connect es error", err)
		return nil
	}
	esClient = client
	return
}

func SendToES(topic string, data []byte) (err error) {
	msg := &LogMessage{}
	msg.Topic = topic
	msg.Message = string(data)

	//首先检查topic索引是否存在
	exists, err := esClient.IndexExists(topic).Do(context.Background())
	if err != nil {
		global.Log.Info(fmt.Sprintf("%s", err))
	}
	if !exists {
		_, err = esClient.CreateIndex(topic).BodyString(mapping).Do(context.Background())
		if err != nil {
			global.Log.Warn(fmt.Sprintf("%s", err))
			return
		}
	}
	_, err = esClient.Index().
		Index(topic).
		BodyJson(msg).
		Do(context.Background())
	if err != nil {
		panic(err)
		return
	}
	return
}

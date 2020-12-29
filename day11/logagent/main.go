package main

import (
	"fmt"
	"time"

	"oldbay.study.com/day11/logagent/conf"

	"gopkg.in/ini.v1"
	"oldbay.study.com/day11/logagent/kafka"
	"oldbay.study.com/day11/logagent/taillog"
)

var (
	cfg = new(conf.AppConf)
)

func run() {
	for {
		select {
		case line := <-taillog.ReadChan():
			kafka.SendToKafka(cfg.KafkaConf.Topic, line.Text)
		default:
			time.Sleep(time.Second)
		}
	}
}

func main() {
	// 0.加载配置文件
	// cfg, err := ini.Load("./conf/config.ini")

	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("init conf failed,err:%v\n", err)
		return
	}
	fmt.Println("加载配置成功")
	// 1.初始化kafka连接
	err = kafka.Init([]string{"172.20.163.200:9092"})
	if err != nil {
		fmt.Printf("init kafka failed,err:%v\n", err)
		return
	}
	fmt.Println("初始化成功")
	// 2.收集日志
	err = taillog.Init(cfg.TaillogConf.FileName)
	if err != nil {
		fmt.Printf("init taillog failed,err:%v\n", err)
		return
	}
	fmt.Println("tail初始化成功")
	run()
}

package main

import (
	"fmt"
	"sync"
	"time"

	"oldbay.study.com/day12/logagent/taillog"
	"oldbay.study.com/day12/logagent/utils"

	"oldbay.study.com/day12/logagent/conf"
	"oldbay.study.com/day12/logagent/etcd"

	"gopkg.in/ini.v1"
	"oldbay.study.com/day12/logagent/kafka"
)

var (
	cfg = new(conf.AppConf)
)

// func run() {
// 	for {
// 		select {
// 		case line := <-taillog.ReadChan():
// 			kafka.SendToKafka(cfg.KafkaConf.Topic, line.Text)
// 		default:
// 			time.Sleep(time.Second)
// 		}
// 	}
// }

func main() {
	// 0.加载配置文件
	// cfg, err := ini.Load("./conf/config.ini")

	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("init conf failed,err:%v\n", err)
		return
	}
	fmt.Println("conf初始化成功")

	// 1.初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.ChanMaxSize)
	if err != nil {
		fmt.Printf("init kafka failed,err:%v\n", err)
		return
	}
	fmt.Println("kafka初始化成功")

	// 2.初始化etcd
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed,err:%v\n", err)
		return
	}
	fmt.Println("etcd初始化成功")

	ipStr, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key, ipStr)
	fmt.Println("==========", etcdConfKey)
	// 2.1 etcd获取配置信息
	logEntryConf, err := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Printf("etcd.GetConf failed,err:%v\n", err)
		return
	}
	fmt.Printf("etcd.GetConf success,%v\n", logEntryConf)

	for index, value := range logEntryConf {
		fmt.Printf("index:%v\n value:%v\n", index, value)
	}

	// 2.2 监视配置变化

	// 收集日志发往kafka
	taillog.Init(logEntryConf)
	newConfChan := taillog.NewConfChan() //从taillog获取通道
	var wg sync.WaitGroup
	wg.Add(1)
	// 派一个哨兵 一直监视着
	fmt.Println(cfg.EtcdConf.Key)
	go etcd.WatchConf(etcdConfKey, newConfChan)
	wg.Wait()
	// // 2.收集日志
	// err = taillog.Init(cfg.TaillogConf.FileName)
	// if err != nil {
	// 	fmt.Printf("init taillog failed,err:%v\n", err)
	// 	return
	// }
	// fmt.Println("tail初始化成功")
	// run()
}

package main

import (
	"fmt"

	"oldbay.study.com/day13/log_transfer/es"

	"oldbay.study.com/day13/log_transfer/kafka"

	"gopkg.in/ini.v1"
	"oldbay.study.com/day13/log_transfer/conf"
)

func main() {

	// 0.加载配置文件
	var cfg = new(conf.LogTransfer)
	err := ini.MapTo(cfg, "./conf/cfg.ini")
	if err != nil {
		fmt.Printf("init config failed, err:%v\n", err)
		return
	}
	fmt.Printf("cfg:%v\n", cfg)
	// 1.初始化ES
	// 1.1 初始化es的连接
	// 1.2 对外提供一个写入ES的函数
	err = es.Init(cfg.ESCfg.Address, cfg.ESCfg.ChanSize, cfg.ESCfg.Nums)
	if err != nil {
		fmt.Printf("init es client failed,err: %v\n", err)
		return
	}
	fmt.Println("init es success.")
	// 2.初始化kafka
	// 2.1 连接kafka, 创建分区的消费者
	// 2.2 每个分区的消费者取出数据发给es
	err = kafka.Init([]string{cfg.KafkaCfg.Address}, cfg.KafkaCfg.Topic)
	if err != nil {
		fmt.Printf("ini kafka consumer failed,err:%v\n", err)
		return
	}
	select {}
}

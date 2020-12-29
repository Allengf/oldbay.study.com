package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// etcd client put/get demo
// use etcd/clientv3
func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"172.20.163.200:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")

	defer cli.Close()

	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	value := `[{"path":"/Users/gaofeng/Data/code/goproject/src/oldbay.study.com/day12/logagent/nginx.log","topic":"web_log"},{"path":"/Users/gaofeng/Data/code/goproject/src/oldbay.study.com/day12/logagent/redis.log","topic":"redis_log"},{"path":"/Users/gaofeng/Data/code/goproject/src/oldbay.study.com/day12/logagent/redis.log","topic":"redis_log"}]`

	_, err = cli.Put(ctx, "/logagent/collect_config", value)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "/logagent/collect_config")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}

package es

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
)

// LogData ...
type LogData struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

//初始化ES,准备接收kafka发来的数据

var (
	client *elastic.Client
	ch     chan *LogData
)

// Init ...
func Init(address string, chanSize, nums int) (err error) {
	if !strings.HasPrefix(address, "httP://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		// Handle error
		return err
	}
	fmt.Println("connect to es success")
	ch = make(chan *LogData, chanSize)
	for i := 0; i < nums; i++ {
		go sendToES()
	}
	// go sendToES()
	return
}

// SendToES 发送数据到ES
func SendToESChan(msg *LogData) {

	ch <- msg
}

func sendToES() {

	for {
		select {
		case msg := <-ch:
			fmt.Println(msg.Topic, msg.Data)
			put1, err := client.Index().Index(msg.Topic).Type("xxx").BodyJson(msg).Do(context.Background())
			if err != nil {
				// Handle error
				fmt.Println(err)
				continue
			}
			fmt.Printf("Indexed Id %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Second)
		}
	}
	// 链式操作

}

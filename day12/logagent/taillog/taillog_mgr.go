package taillog

import (
	"fmt"
	"time"

	"oldbay.study.com/day12/logagent/etcd"
)

var tskMgr *tailLogMgr

type tailLogMgr struct {
	logEntry    []*etcd.LogEntry
	tskMap      map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(LogEntryConf []*etcd.LogEntry) {
	tskMgr = &tailLogMgr{
		logEntry:    LogEntryConf,
		tskMap:      make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry),
	}
	for _, logEntry := range LogEntryConf {

		tailObj := NewTailTask(logEntry.Path, logEntry.Topic)
		mk := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)
		fmt.Println("++++++", mk)
		tskMgr.tskMap[mk] = tailObj
	}
	go tskMgr.run()
}

func (t *tailLogMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				fmt.Println("----------", mk)
				_, ok := t.tskMap[mk]
				if ok {
					continue
				} else {
					tailObj := NewTailTask(conf.Path, conf.Topic)
					t.tskMap[mk] = tailObj
				}
			}
			for _, c1 := range t.logEntry {
				isDelete := true
				for _, c2 := range newConf {
					if c2.Path == c1.Path && c2.Topic == c1.Topic {
						isDelete = false
						continue
					}
				}
				if isDelete {
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					t.tskMap[mk].cancelFunc()
				}
			}

			fmt.Println("新的配置来了！", newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}

// func PushNewConf() chan<- []*etcd.LogEntry {
// 	return tskMgr.newConfChan
// }

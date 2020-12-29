package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func f2(ctx context.Context) {
	defer wg.Done()
LOOP:
	for {
		fmt.Println("显示")
		time.Sleep(time.Millisecond * 600)
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
	}
}

func f(ctx context.Context) {
	defer wg.Done()
	go f2(ctx)
LOOP:
	for {
		fmt.Println("刘")
		time.Sleep(time.Millisecond * 500)
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go f(ctx)
	time.Sleep(time.Second * 5)
	cancel()
	wg.Wait()
}

package main

import (
	"cron"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

//corn重要的api
//cron.Every(p time.Duration) Schedule  通过传入一个时间段，获得了一个实现了next方法的接口(计算任务结束的时间的方法)
//cron.AddFunc(schedule Schedule, j func())   将cron.Every(p time.Duration)获得schedule，和定时的任务函数传入
//时间精确度到s

//模拟最简单的打印
//使用wg为了防止主函数过早退出，导致g没有成功运行
func main() {
	wg.Add(1)
	c := cron.New()
	c.AddFunc(cron.Every(5*time.Second), func() {
		fmt.Println("runs every 5 seconds.")
	})
	c.Start()

	wg.Wait()
}

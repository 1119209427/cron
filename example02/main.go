package main

import (
	"cron"
	"fmt"
	"sync"
	"time"
)

//自定义任务
//实现job接口的类型即可
/*type Job interface {
	Run()
}*/

type Hello struct {
	Name string
}

func (h Hello) Run() {
	fmt.Printf("hello %s\n", h.Name)
}

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	c := cron.New()
	g1 := Hello{Name: "jc"}
	g2 := Hello{Name: "ke qin"}
	c.AddTask(cron.Every(5*time.Second), g1)
	c.AddTask(cron.Every(3*time.Second), g2)
	c.Start()
	wg.Wait()
}

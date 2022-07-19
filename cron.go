package cron

import (
	"sort"
	"time"
)

// Task 每个任务的结构体
type Task struct {
	Job      Job       //任务内容
	Schedule Schedule  //用来计算下次执行时间
	Prev     time.Time //任务上次执行的时间
	Next     time.Time // 任务下次执行的时间
}

type Cron struct {
	Tasks   []*Task       //要执行的任务
	Running bool          //是否在执行任务
	add     chan *Task    //向cron加入任务的管道
	stop    chan struct{} //结束cron的管道  因为空结构体不占据内存空间，因此被广泛作为各种场景下的占位符使用。一是节省资源，二是空结构体本身就具备很强的语义，即这里不需要任何值，仅作为占位符。
}

//根据注释可知，在等待给定的一段时间后，向返回值发送当前时间，返回值是一个单向只读通道
//传入每一个任务的周期，获取任务结束的时间
var after = time.After

//为了使用sort接口
type byTime []*Task

func (b byTime) Len() int      { return len(b) }
func (b byTime) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byTime) Less(i, j int) bool {

	if b[i].Next.IsZero() {
		return false
	}
	if b[j].Next.IsZero() {
		return true
	}

	return b[i].Next.Before(b[j].Next)
}
func New() *Cron {
	return &Cron{
		add:  make(chan *Task),
		stop: make(chan struct{}),
	}

}
func (c *Cron) Start() {
	c.Running = true
	go c.run()
}
func (c *Cron) Stop() {
	//没有任务运行
	if !c.Running {
		return
	}
	//有任务运行，通过管道通信来结束
	c.stop <- struct{}{}
}
func (c *Cron) AddTask(schedule Schedule, job Job) {

	task := &Task{Schedule: schedule, Job: job}
	if !c.Running {
		c.Tasks = append(c.Tasks, task)
		return
	}
	c.add <- task

}
func (c *Cron) AddFunc(schedule Schedule, j func()) {
	c.AddTask(schedule, JobFunc(j))
}

//核心方法，实现cron的运行
func (c *Cron) run() {
	var effective time.Time   //记录第一个任务的next
	now := time.Now().Local() //活得现在的时间
	for _, task := range c.Tasks {
		//计算出每个任务的next
		task.Next = task.Schedule.Next(now)
	}
	//对 tasks切片进行排序
	for {
		sort.Sort(byTime(c.Tasks))
		//获取tasks切片中的第一个对象就是要执行的第一个任务
		if len(c.Tasks) > 0 {
			effective = c.Tasks[0].Next
		} else {
			effective = now.AddDate(15, 0, 0) //防止浪费内存
		}
		select {
		case now = <-after(effective.Sub(now)): //effective.Sub(now) 返回当前时间到执行第一个任务的时间差
			for _, task := range c.Tasks {
				if task.Next != effective {
					//说明没到这个任务执行
					break
				}
				task.Prev = now
				task.Next = task.Schedule.Next(now)
				//运行定时任务
				go task.Job.Run()

			}
		case task := <-c.add:
			task.Next = task.Schedule.Next(now)
			//将新任务加入
			c.Tasks = append(c.Tasks, task)
		case <-c.stop:
			break

		}
	}

}

// Job 抽象定义  方便使用者扩展
type Job interface {
	Run()
}

// JobFunc 封装一个函数类型的接口，方便外部使用
type JobFunc func()

func (j JobFunc) Run() {
	j()
}

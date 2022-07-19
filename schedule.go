package cron

import "time"

type Schedule interface {
	Next(t time.Time) time.Time //根据传入的时间段，计算任务的下一次执行时间
}
type periodicSchedule struct {
	period time.Duration
}

func Every(p time.Duration) Schedule {
	p = p - time.Duration(p.Nanoseconds())%time.Second // truncates up to seconds

	return &periodicSchedule{
		period: p,
	}
}
func (ps periodicSchedule) Next(t time.Time) time.Time {
	return t.Truncate(time.Second).Add(ps.period)
}

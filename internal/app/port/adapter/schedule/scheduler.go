package schedule

import "github.com/robfig/cron/v3"

type Task interface {
	cron.Job
}

type Scheduler interface {
	AddTask(task Task)
}

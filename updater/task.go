package updater

import (
	"fmt"
	"sync"
	"time"
)

type LogType int8

const (
	_ LogType = iota
	LtDebug
	LtInfo
	LtWarning
	LtError
)

type TaskLog struct {
	LogType LogType
	Message string
}

type RunnableTask interface {
	Run()
	GetName() string
}

type Task struct {
	Name string

	mux         sync.Mutex
	waitCond    *sync.Cond
	concurrency int
	delay       int
	waitGroup   sync.WaitGroup
	queue       chan int
	suspended   bool

	start    *int64
	end      *int64
	progress *int64
	logChan  *chan TaskLog

	rateLimiter chan int
}

func (task *Task) GetName() string {
	return task.Name
}

func (task *Task) rateLimitWorker() {
	for range task.rateLimiter {
		time.Sleep(time.Duration(task.delay) * time.Millisecond)
	}
}

func (task *Task) log(logType LogType, message string, args ...interface{}) {
	if task.logChan != nil {
		*task.logChan <- TaskLog{
			LogType: logType,
			Message: fmt.Sprintf(message, args...),
		}
	}
}

func (task *Task) resume() {
	time.Sleep(15 * time.Minute) // Wait for an hour
	task.log(LtWarning, "RESUMING ALL WORKERS\n")
	task.suspended = false
	task.waitCond.Broadcast()
}

func (task *Task) suspend(workerId int) {
	task.mux.Lock()
	defer task.mux.Unlock()

	if task.suspended {
		task.log(LtDebug, "[Worker %d]ALREADY SUSPENDED\n", workerId)
		return
	}

	task.log(LtWarning, "[Worker %d]SUSPENDING ALL WORKERS\n", workerId)
	task.suspended = true

	go task.resume()
}

package updater

import (
	"os"
	"sync"
	"sync/atomic"
)

type TaskManager struct {
	Concurrency int
	Delay       int
	LogLevel    LogType
	CurrentTask RunnableTask
	Progress    int
	Status      bool

	taskList []RunnableTask

	taskStart    int64
	taskEnd      int64
	taskProgress int64
	logChannel   *chan TaskLog

	CachedRequests   int32
	UncachedRequests int32
	FailedRequests   int32
}

func NewTaskManager(concurrency int, delay int, logLevel LogType) *TaskManager {
	return &TaskManager{
		Concurrency:      concurrency,
		Delay:            delay,
		LogLevel:         logLevel,
		Progress:         0,
		Status:           true,
		CachedRequests:   0,
		UncachedRequests: 0,
		FailedRequests:   0,
	}
}

func (manager *TaskManager) incCachedRequests() {
	atomic.AddInt32(&manager.CachedRequests, 1)
}

func (manager *TaskManager) incUncachedRequests() {
	atomic.AddInt32(&manager.UncachedRequests, 1)
}

func (manager *TaskManager) incFailedRequests() {
	atomic.AddInt32(&manager.FailedRequests, 1)
}

func (manager *TaskManager) AddIndexTask(name string, indexMethod string, indexCollection string, itemMethod string, updateCallback ItemCallback) {
	task := &IndexTask{
		Task: Task{
			Name:        name,
			mux:         sync.Mutex{},
			concurrency: manager.Concurrency,
			delay:       manager.Delay,

			logChan: manager.logChannel,
			manager: manager,
		},
		IndexMethod:     indexMethod,
		IndexCollection: indexCollection,
		ItemMethod:      itemMethod,
		ItemCallback:    updateCallback,
	}
	manager.taskList = append(manager.taskList, task)
}

func (manager *TaskManager) AddIndexTaskLimited(name string, indexMethod string, indexCollection string, itemMethod string, updateCallback ItemCallback, concurrency int) {
	task := &IndexTask{
		Task: Task{
			Name:        name,
			mux:         sync.Mutex{},
			concurrency: concurrency,
			delay:       manager.Delay,

			logChan: manager.logChannel,
			manager: manager,
		},
		IndexMethod:     indexMethod,
		IndexCollection: indexCollection,
		ItemMethod:      itemMethod,
		ItemCallback:    updateCallback,
	}
	manager.taskList = append(manager.taskList, task)
}

func (manager *TaskManager) AddSearchTask(name string, indexMethod string, itemMethod string, updateCallback ItemCallback) {
	task := &SearchTask{
		Task: Task{
			Name:        name,
			mux:         sync.Mutex{},
			concurrency: manager.Concurrency,
			delay:       manager.Delay,

			logChan: manager.logChannel,
			manager: manager,
		},
		SearchMethod: indexMethod,
		ItemMethod:   itemMethod,
		ItemCallback: updateCallback,
	}
	manager.taskList = append(manager.taskList, task)
}

func (manager *TaskManager) AddSimpleTask(name string, method SimpleMethod) {
	task := &SimpleTask{
		Name:   name,
		Method: method,
	}
	manager.taskList = append(manager.taskList, task)
}

func (manager *TaskManager) Run() {
	maxProgress := len(manager.taskList)
	for progress, task := range manager.taskList {
		manager.Progress = progress * 100.0 / maxProgress
		manager.CurrentTask = task
		task.Run()
	}
	os.Exit(0)
}

package updater

import (
	"fmt"
	blizzard_api "github.com/francis-schiavo/blizzard-api-go"
	"log"
	"reflect"
	"sync"
	"wow-query-updater/connections"
	"wow-query-updater/datasets"
)

type SearchTask struct {
	Task
	SearchMethod string
	ItemMethod   string
	ItemCallback ItemCallback
}

func (task *SearchTask) worker(workerId int) {
	endpointInterface := reflect.ValueOf(connections.WowClient).MethodByName(task.ItemMethod)

	for id := range task.queue {
		args := []reflect.Value{
			reflect.ValueOf(id),
			reflect.ValueOf((*blizzard_api.RequestOptions)(nil)),
		}
		task.log(LT_DEBUG, "[Worker %d] Processing %s %d\n", workerId, task.Name, id)

		response := endpointInterface.Call(args)[0].Interface().(*blizzard_api.ApiResponse)
		if !response.Cached {
			task.rateLimiter <- 1
		}

		task.log(LT_DEBUG, "[Worker %d] Finished processing %s %d\n", workerId, task.Name, id)
		if response.Status == 200 {
			task.ItemCallback(response)
			task.log(LT_INFO, "Updated %s %d successfully!\n", task.Name, id)
			task.waitGroup.Done()
		} else if response.Status == 429 {
			// Insert the failed id into the queue to retry later
			task.queue <- id
			// Suspend all goroutines temporarily
			task.suspend(workerId)
		} else {
			task.log(LT_INFO, "Failed to update %s %d with status: %d\n", task.Name, id, response.Status)
			task.waitGroup.Done()
		}

		// Wait for a while after a "too many requests" response
		if task.suspended {
			// Wait until it is resumed
			task.log(LT_DEBUG, "[Worker %d] Waiting\n", workerId)
			task.waitCond.L.Lock()
			task.waitCond.Wait()
			task.waitCond.L.Unlock()
			task.log(LT_DEBUG, "[Worker %d] Resumed\n", workerId)
		}
	}
	task.log(LT_DEBUG, "[Worker %d] Exiting\n", workerId)
}

func (task *SearchTask) Run() {
	task.waitGroup = sync.WaitGroup{}
	task.queue = make(chan int)
	m := &sync.Mutex{}
	task.waitCond = sync.NewCond(m)
	task.suspended = false

	endpointInterface := reflect.ValueOf(connections.WowClient).MethodByName(task.SearchMethod)

	for w := 1; w <= task.concurrency-1; w++ {
		go task.worker(w)
	}

	task.rateLimiter = make(chan int, task.concurrency-1)
	go task.rateLimitWorker()

	var jsonData datasets.SearchResult
	lastID := 0
	for {
		args := []reflect.Value{
			reflect.ValueOf(&blizzard_api.RequestOptions{
				QueryString: map[string]string {
					"orderby": "id",
					"id": fmt.Sprintf("[%d,]", lastID + 1),
					"_pageSize": "1000",
				},
			}),
		}
		response := endpointInterface.Call(args)[0].Interface().(*blizzard_api.ApiResponse)

		if response.Status != 200 {
			log.Fatalf("Failed to obtain search data with status: %d.\n", response.Status)
		}

		response.Parse(&jsonData)

		if len(jsonData.Results) == 0 {
			fmt.Sprintln("No more data")
			break
		}

		for _, item := range jsonData.Results {
			task.waitGroup.Add(1)
			lastID = item.Data.ID
			task.queue <- lastID
		}
	}

	task.waitGroup.Wait()
	task.log(LT_INFO, "Task %s finished.", task.Name)
	close(task.queue)
	close(task.rateLimiter)
}

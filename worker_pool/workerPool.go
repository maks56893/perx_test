package worker_pool

import (
	"fmt"
	"perx_test/cache"
	"perx_test/consts"
	"perx_test/structs"
	"strconv"
	"sync"
	"time"
)

// WorkerPool - структура пула воркеров
type WorkerPool struct {
	maxWorkers     int
	maxQuerySize   int
	workerQuit     chan struct{}
	updateChannel  chan workerTaskStruct
	wg             sync.WaitGroup
	workerWG       sync.WaitGroup
	completedTasks *cache.Cache
}

type workerTaskStruct struct {
	workerFunc func(*structs.Task)
	userData   *structs.Task
}

var singleWorkerPool *WorkerPool

// NewWorkerPool - Создание объекта пула воркеров 
func NewWorkerPool(cache *cache.Cache) *WorkerPool {
	if singleWorkerPool == nil {
		singleWorkerPool = &WorkerPool{
			updateChannel: make(chan workerTaskStruct),
			maxWorkers:    consts.MaxWorkers,
			maxQuerySize:  consts.MaxQuerySize,
			completedTasks: cache,
		}
	}
	return singleWorkerPool
}

// Push - добавить задачу в очередь
func (s *WorkerPool) Push(f func(*structs.Task), task *structs.Task) {
	s.updateChannel <- workerTaskStruct{
		workerFunc: f,
		userData:   task,
	}
}

func (s *WorkerPool) worker(ch chan workerTaskStruct) {

	for {
		select {
		case <-s.workerQuit:
			s.workerWG.Done()
			return
		case ev, ok := <-ch:
			if ok {
				ev.workerFunc(ev.userData)
				s.completedTasks.Set(strconv.Itoa(ev.userData.TaskNumber), ev.userData, time.Duration(ev.userData.Ttl) * time.Second)
			} else {
				return
			}
		}
	}
}

// Start - запустить пул воркеров
func (s *WorkerPool) Start() {

	s.updateChannel = make(chan workerTaskStruct, s.maxQuerySize)
	s.workerQuit = make(chan struct{})
	s.workerWG.Add(s.maxWorkers)

	for i := 0; i < s.maxWorkers; i++ {
		go s.worker(s.updateChannel)
	}

}

// Stop - остановить пул воркеров
func (s *WorkerPool) Stop() {
	fmt.Println("Остановка пула воркеров...")
	close(s.workerQuit)

	s.wg.Wait()
}

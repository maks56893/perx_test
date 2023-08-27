package task_controller

import (
	//"container/list"
	"context"
	"fmt"
	"perx_test/cache"
	"perx_test/structs"
	"perx_test/worker_pool"
	"strconv"
)

// TaskCtrl - глобальный объект контроллера задач. Должен инициализироваться при запуске
var TaskCtrl *TaskController

type TaskController struct {
	ctx         context.Context
	taskCounter int
	Tasks       *cache.Cache
	workerPool  *worker_pool.WorkerPool
}

// InitTaskController - если ранее не был создан контроллер задач, то возвращает новый объект, иначе возвращает уже созданный
func InitTaskController(ctx context.Context, cache *cache.Cache) *TaskController {
	if TaskCtrl == nil {
		workerPool := worker_pool.NewWorkerPool(cache)
		workerPool.Start()

		TaskCtrl = &TaskController{
			ctx:        ctx,
			workerPool: workerPool,
			Tasks:      cache,
		}
	}
	return TaskCtrl
}

// Stop - остановка контроллера выполнения задач и завершения текущих задач
func (tc *TaskController) Stop() {
	tc.workerPool.Stop()
}

// Append - поставить задачу в очередь на выполнение
func (tc *TaskController) Append(task *structs.Task) {
	select {
	case <-tc.ctx.Done():
		fmt.Println("Идет завершение работы, новые задачи не будут обработаны")
		return
	default:
		tc.taskCounter++
		task.TaskNumber = tc.taskCounter
		tc.Tasks.Set(strconv.Itoa(tc.taskCounter), task, 0) // добавляем в кеш чтобы отслеживать не начатые таски, такие записи не истекают
		tc.workerPool.Push(CountSequense, task)
	}
}

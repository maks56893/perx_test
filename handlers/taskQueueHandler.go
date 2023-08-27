package handlers

import (
	"net/http"
	"perx_test/structs"
	"perx_test/task_controller"
)

type tasksQueueResp struct {
	Queue []*structs.Task `json:"tasks"`
}

// WriteError - записывает ошибку в структуру ответа. В этом хендлере нет может быть ошибки
func (resp *tasksQueueResp) WriteError(err error) {}

// GetTasksHandler - хендлер /stat. Тела нет
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := make([]*structs.Task, 0)
	task_controller.TaskCtrl.Tasks.Range(func(key string, value interface{}) bool {
		tasks = append(tasks, value.((*structs.Task)))
		return true
	})

	resp := &tasksQueueResp{
		Queue: tasks,
	}
	registerSuccess(w, resp)
}

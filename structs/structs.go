package structs

import "time"

type taskStatus string

func (status taskStatus) ToInt() int {
	switch status {
	case InQueue:
		return 1
	case InProgres:
		return 2
	case Completed:
		return 3
	default:
		return 0
	}
}

// Статусы выполнения задач
const (
	InQueue   taskStatus = "in_queue"
	InProgres taskStatus = "in_progress"
	Completed taskStatus = "completed"
)

// Task - структура задачи арифметической прогрессии
type Task struct {
	N   int     `json:"n"`   // количество элементов
	D   float64 `json:"d"`   // дельта между элементами последовательности
	N1  float64 `json:"n1"`  // Стартовое значение
	L   float64 `json:"l"`   // интервал в секундах между итерациями
	Ttl float64 `json:"ttl"` // время хранения результата в секундах

	TaskNumber       int        `json:"task_number"`
	Status           taskStatus `json:"status"`
	CurrentIteration int        `json:"current_iteration"`
	GetTaskTime      time.Time  `json:"get_task_time"`
	StartTime        time.Time  `json:"start_time,omitempty"`
	EndTime          time.Time  `json:"end_time,omitempty"`

	Result []float64 `json:"result,omitempty"`
}

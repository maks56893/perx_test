package task_controller

import (
	"perx_test/structs"
	"time"
)

// CountSequense - функция-обработчик задачи на подсчет арифметической прогрессии. Внутри, кроме результата, проставляет информацию о задаче
func CountSequense(task *structs.Task) {
	task.Status = structs.InProgres
	task.StartTime = time.Now()

	result := make([]float64, 0, task.N)
	result = append(result, task.N1)
	ticker := time.NewTicker(time.Duration(task.L) * time.Second)
	for i := 1; i < task.N; i++ {
		<-ticker.C
		val := result[i-1] + task.D
		result = append(result, val)
		task.CurrentIteration = i+1
	}

	task.EndTime = time.Now()
	task.Status = structs.Completed
	task.Result = result
}

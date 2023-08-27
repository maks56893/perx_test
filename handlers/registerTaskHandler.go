package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"perx_test/task_controller"
	"perx_test/structs"
	"time"
)

type registerTaskResp struct {
	ErrorString string `json:"error_string"`
}

func (resp *registerTaskResp) WriteError(err error) {
	resp.ErrorString = err.Error()
}

type registerTaskReq struct {
	N      int     `json:"n"`   // количество элементов
	D      float64 `json:"d"`   // дельта между элементами последовательности
	N1     float64 `json:"n1"`  // Стартовое значение
	L      float64 `json:"l"`   // интервал в секундах между итерациями
	Ttl    float64 `json:"ttl"` // время хранения результата в секундах
}

// RegisterTaskHandler - функция-хендлер /register
// {
//     "n": 1,
//     "d": 2,
//     "n1": 3,
//     "l": 4,
//     "ttl": 5
// }
func RegisterTaskHandler(w http.ResponseWriter, r *http.Request) {
	resp := new(registerTaskResp)

	err := r.ParseForm()
	if err != nil {
		registerError(err, w, resp)
		return
	}
	
	bodyRaw := r.Body
	buf := new(bytes.Buffer)
	buf.ReadFrom(bodyRaw)
	reqBytes := buf.Bytes()

	req := registerTaskReq{}
	err = json.Unmarshal(reqBytes, &req)
	if err != nil {
		registerError(err, w, resp)
		return
	}

	task := structs.Task{
		N: req.N,
		D: req.D,
		N1: req.N1,
		L: req.L,
		Ttl: req.Ttl,
		Status: structs.InQueue,
		GetTaskTime: time.Now(),
	}
	task_controller.TaskCtrl.Append(&task)
	registerSuccess(w, resp)
}
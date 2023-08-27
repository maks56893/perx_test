package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"syscall"

	"perx_test/cache"
	"perx_test/consts"
	"perx_test/handlers"
	"perx_test/task_controller"
)

type Application struct {
	ctx            context.Context
	cancelFunc     context.CancelFunc
	taskController *task_controller.TaskController
	server         *http.Server
}

// NewApp - создание нового объекта приложения
func NewApp() *Application {
	if len(os.Args) != 2 {
		err := fmt.Errorf("Количество аргументов запуска должно быть равно двум")
		log.Fatal(err)
	}
	maxWorkers, err := strconv.Atoi(os.Args[1])
	if err != nil {
		err = fmt.Errorf("Не удалось конвертировать аргумент командной строки %v в int значение: %v", os.Args[1], err)
		log.Fatal(err)
	}
	consts.MaxWorkers = maxWorkers

	ctx, cancelFunc := context.WithCancel(context.Background())
	cacheForCompletedTasks := cache.New(consts.DefaultExpiration, consts.CleanupInterval)
	taskCtrl := task_controller.InitTaskController(ctx, cacheForCompletedTasks)
	srv := &http.Server{Addr: ":9000"}

	return &Application{
		ctx:            ctx,
		cancelFunc:     cancelFunc,
		taskController: taskCtrl,
		server:         srv,
	}
}

// Stop - остановка приложения
func (app *Application) Stop(ch chan os.Signal) {
	<-ch
	fmt.Println("Завершение работы...")
	app.server.Shutdown(app.ctx)
	app.cancelFunc()
	app.taskController.Stop()
}

func main() {
	app := NewApp()
	// Канал для ожидания сигнала Ctrl+C
	ch := make(chan os.Signal, 1)

	http.HandleFunc("/register", handlers.RegisterTaskHandler)
	http.HandleFunc("/stat", handlers.GetTasksHandler)
	go func(ch chan os.Signal) {
		err := app.server.ListenAndServe() // задаем слушать порт
		if err != http.ErrServerClosed {
			log.Fatal("ListenAndServe: ", err)
		}
	}(ch)

	fmt.Println("Ожидание завершения")

	// Завершение Ctrl+C
	signal.Notify(ch, syscall.SIGINT)

	app.Stop(ch)
}

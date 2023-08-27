package consts

import "time"

const (
	DefaultExpiration = 0 * time.Second // DefaultExpiration - время протухания объектов в кеше по умолчанию. Если равно 0, то объекты не удаляются
	CleanupInterval = 60 * time.Second // CleanupInterval - интервал времени, через который запустится сборщик протухших значений в кеше
	MaxQuerySize = 255 // MaxQuerySize - максимальное количество элементов в очереди в пуле воркеров
)

var (
	MaxWorkers int // MaxWorkers - максимальное количество воркеров. Передается аргументом при запуске, если не задано, то выдаем ошибку и завершаем работу
)
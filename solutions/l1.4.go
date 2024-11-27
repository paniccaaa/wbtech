package solutions

import (
	"flag"
)

// Реализовать постоянную запись данных в канал (главный поток).
// Реализовать набор из N воркеров, которые читают произвольные данные из канала и выводят в stdout.
// Необходима возможность выбора количества воркеров при старте.
// Программа должна завершаться по нажатию Ctrl+C.
// Выбрать и обосновать способ завершения работы всех воркеров

var numWorkers = flag.Int("worker", 3, "num of workers")

func Solve4() {
	flag.Parse()
	mainChannel := make(chan int)

	for {
		mainChannel <- 1
	}

}

func worker() {

}

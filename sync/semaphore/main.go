package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	MaxResources = 3
	NumWorkers   = 10
)

func main() {
	var wg sync.WaitGroup

	// Создаем семафор с ограничением MaxResources
	sem := semaphore.NewWeighted(int64(MaxResources))

	wg.Add(NumWorkers)

	// Запускаем воркеров
	for i := 0; i < NumWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			workerSem(workerID, sem)
		}(i)
	}

	wg.Wait()
	fmt.Println("\nВсе воркеры завершили работу")
}

func workerSem(id int, sem *semaphore.Weighted) {
	// Пытаемся получить ресурс (1 единицу веса)
	fmt.Printf("Worker %d: пытаюсь получить ресурс\n", id)

	if err := sem.Acquire(nil, 1); err != nil {
		fmt.Printf("Worker %d: ошибка получения ресурса: %v\n", id, err)
		return
	}

	fmt.Printf("Worker %d: получил ресурс\n", id)

	// Выполняем работу
	workTime := time.Duration(500+id*100) * time.Millisecond
	fmt.Printf("Worker %d: работаю %v\n", id, workTime)
	time.Sleep(workTime)

	// Освобождаем ресурс
	sem.Release(1)
	fmt.Printf("Worker %d: освободил ресурс\n", id)
}

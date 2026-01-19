package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	MaxResources = 3
	NumWorkers   = 10
)

func main() {
	var wg sync.WaitGroup

	// Создаем пул ресурсов с использованием каналов
	resourcePool := make(chan struct{}, MaxResources)

	// Заполняем пул ресурсами
	for i := 0; i < MaxResources; i++ {
		resourcePool <- struct{}{}
	}

	wg.Add(NumWorkers)

	// Запускаем воркеров
	for i := 0; i < NumWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			worker(workerID, resourcePool)
		}(i)
	}

	wg.Wait()
	fmt.Println("\nВсе воркеры завершили работу")
}

func worker(id int, pool chan struct{}) {
	// Пытаемся получить ресурс из пула
	fmt.Printf("Worker %d: пытаюсь получить ресурс\n", id)

	// Блокируемся, пока не получим ресурс
	<-pool
	fmt.Printf("Worker %d: получил ресурс (свободно: %d/%d)\n",
		id, len(pool), cap(pool))

	// Выполняем работу
	workTime := time.Duration(500+id*100) * time.Millisecond
	fmt.Printf("Worker %d: работаю %v\n", id, workTime)
	time.Sleep(workTime)

	// Возвращаем ресурс в пул
	pool <- struct{}{}
	fmt.Printf("Worker %d: вернул ресурс (свободно: %d/%d)\n",
		id, len(pool), cap(pool))
}

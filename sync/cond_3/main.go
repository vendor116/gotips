package main

import (
	"fmt"
	"sync"
	"time"
)

const MaxResources = 3
const NumWorkers = 10

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	// Создаем провайдера ресурсов
	resourceProvider := NewResourceProvider(MaxResources, cond)

	wg.Add(NumWorkers)

	// Запускаем воркеров
	for i := 0; i < NumWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			worker := NewWorker(workerID, resourceProvider)
			worker.Run()
		}(i)
	}

	wg.Wait()
	fmt.Println("\nВсе воркеры завершили работу")
}

type ResourceProvider struct {
	maxResources       int
	availableResources int
	cond               *sync.Cond
}

func NewResourceProvider(maxResources int, cond *sync.Cond) *ResourceProvider {
	return &ResourceProvider{
		maxResources:       maxResources,
		availableResources: maxResources,
		cond:               cond,
	}
}

func (rp *ResourceProvider) AvailableResources() int {
	return rp.availableResources
}

func (rp *ResourceProvider) AcquireResource() {
	if rp.availableResources > 0 {
		rp.availableResources--
	}
}

func (rp *ResourceProvider) ReleaseResource() {
	if rp.availableResources < rp.maxResources {
		rp.availableResources++
	}
}

type Worker struct {
	id int
	rp *ResourceProvider
}

func NewWorker(workerID int, rp *ResourceProvider) *Worker {
	return &Worker{
		id: workerID,
		rp: rp,
	}
}

func (w *Worker) Run() {
	// Пытаемся получить ресурс
	w.acquireResource()

	// Выполняем работу
	w.doWork()

	// Освобождаем ресурс
	w.releaseResource()
}

func (w *Worker) acquireResource() {
	w.rp.cond.L.Lock()
	defer w.rp.cond.L.Unlock()

	// Ждем, пока не освободится ресурс
	for w.rp.AvailableResources() == 0 {
		fmt.Printf("Worker %d: жду ресурс (доступно: %d/%d)\n",
			w.id, w.rp.AvailableResources(), w.rp.maxResources)
		w.rp.cond.Wait()
	}

	// Получаем ресурс
	w.rp.AcquireResource()
	fmt.Printf("Worker %d: получил ресурс (осталось: %d/%d)\n",
		w.id, w.rp.AvailableResources(), w.rp.maxResources)
}

func (w *Worker) doWork() {
	// Симулируем работу
	workTime := time.Duration(500+w.id*100) * time.Millisecond
	fmt.Printf("Worker %d: выполняю работу %v\n", w.id, workTime)
	time.Sleep(workTime)
	fmt.Printf("Worker %d: работа завершена\n", w.id)
}

func (w *Worker) releaseResource() {
	w.rp.cond.L.Lock()
	defer w.rp.cond.L.Unlock()

	// Освобождаем ресурс
	w.rp.ReleaseResource()
	fmt.Printf("Worker %d: освободил ресурс (теперь: %d/%d)\n",
		w.id, w.rp.AvailableResources(), w.rp.maxResources)

	// Уведомляем всех ожидающих воркеров
	w.rp.cond.Broadcast() // Используем Broadcast вместо Signal для нескольких воркеров
}

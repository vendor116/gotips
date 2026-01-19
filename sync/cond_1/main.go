package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	// Флаг для проверки условия (рекомендуемый паттерн)
	conditionMet := false

	wg.Add(2)

	// Goroutine 1 - ждет условие
	go func() {
		fmt.Println("Goroutine 1: запущена")
		defer wg.Done()

		mu.Lock()
		fmt.Println("Goroutine 1: ждет выполнения условия")

		// Правильный паттерн: ожидание в цикле с проверкой условия
		for !conditionMet {
			cond.Wait()
		}

		fmt.Println("Goroutine 1: условие выполнено!")
		mu.Unlock()

		fmt.Println("Goroutine 1: завершена")
	}()

	// Goroutine 2 - сигнализирует об условии
	go func() {
		fmt.Println("Goroutine 2: запущена")
		defer wg.Done()

		// Симулируем какую-то работу
		fmt.Println("Goroutine 2: выполняет работу...")
		time.Sleep(2 * time.Second)

		mu.Lock()
		fmt.Println("Goroutine 2: отправляет сигнал")

		// Устанавливаем флаг и отправляем сигнал
		conditionMet = true
		cond.Signal()

		fmt.Println("Goroutine 2: сигнал отправлен")
		mu.Unlock()

		fmt.Println("Goroutine 2: завершена")
	}()

	// Ожидаем завершения всех горутин
	wg.Wait()
	fmt.Println("\nВсе горутины завершены")
}

package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// SimpleGenerator генератор, останавливаемый по контексту, с возвратом канала,
// не блокируется при ожидании читателя
func SimpleGenerator(ctx context.Context, delay time.Duration) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		var h string
		for {
			// если есть контекст, первым делом проверяем ошибку
			if ctx.Err() != nil {
				fmt.Println(fmt.Errorf("context error: %w", ctx.Err()))
				return
			}

			h = randomHost()

			select {
			case <-ctx.Done():
				fmt.Println(fmt.Errorf("context error: %w", ctx.Err()))
				return
			case ch <- h:
				time.Sleep(delay)
			default:
				// не блокируемся, если нет ожидающего читателя
				fmt.Println("empty read-made channel:", h)
				time.Sleep(delay)
			}
		}
	}()

	return ch
}

var hosts = []string{
	"google.com",
	"facebook.com",
	"twitter.com",
	"instagram.com",
}

func randomHost() string {
	return hosts[rand.Intn(len(hosts))]
}

var (
	timeout    = 1 * time.Second
	writeDelay = 10 * time.Millisecond
	readDelay  = 30 * time.Millisecond
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for v := range SimpleGenerator(ctx, writeDelay) {
		fmt.Println("host:", v)
		time.Sleep(readDelay)
	}
}

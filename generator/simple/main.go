package main

import (
	"context"
	"log"
	randv2 "math/rand/v2"
	"time"
)

// HostGenerator генератор, останавливаемый по контексту, с возвратом канала,
// не блокируется при ожидании читателя.
func HostGenerator(ctx context.Context, delay time.Duration, hosts []string) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		var h string
		for {
			// если есть контекст, первым делом проверяем ошибку
			if ctx.Err() != nil {
				log.Printf("context error: %v\n", ctx.Err())
				return
			}

			h = getHost(hosts)

			select {
			case <-ctx.Done():
				log.Printf("context error: %v\n", ctx.Err())
				return
			case ch <- h:
				time.Sleep(delay)
			default:
				// не блокируемся, если нет ожидающего читателя
				log.Printf("empty read-made channel: %v\n", h)
				time.Sleep(delay)
			}
		}
	}()

	return ch
}

func getHost(hosts []string) string {
	return hosts[randv2.IntN(len(hosts))] //nolint:gosec // для простого примера достаточно
}

func main() {
	var (
		hosts = []string{
			"google.com",
			"facebook.com",
			"twitter.com",
			"instagram.com",
		}

		timeout    = time.Second * 1
		writeDelay = time.Millisecond * 10
		readDelay  = time.Millisecond * 30
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for h := range HostGenerator(ctx, writeDelay, hosts) {
		log.Printf("host: %v\n", h)
		time.Sleep(readDelay)
	}
}

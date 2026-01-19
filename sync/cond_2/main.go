package main

import (
	"fmt"
	"sync"
	"time"
)

const MaxMessageChannelSize = 5

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	messageChannel := NewMessageChannel(MaxMessageChannelSize)

	producer := NewProducer(cond, messageChannel)
	consumer := NewConsumer(cond, messageChannel)

	wg.Add(2)

	// Producer goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			producer.Produce(fmt.Sprintf("Message %d", i))
		}
		fmt.Println("Producer: завершил производство 10 сообщений")
	}()

	// Consumer goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			consumer.Consume()
		}
		fmt.Println("Consumer: завершил потребление 10 сообщений")
	}()

	wg.Wait()
	fmt.Println("\nПрограмма завершена. Оставшиеся сообщения в буфере:", len(messageChannel.buffer))
}

type MessageChannel struct {
	maxBufferSize int
	buffer        []string
}

func NewMessageChannel(size int) *MessageChannel {
	return &MessageChannel{
		maxBufferSize: size,
		buffer:        make([]string, 0, size),
	}
}

func (mc *MessageChannel) IsEmpty() bool {
	return len(mc.buffer) == 0
}

func (mc *MessageChannel) IsFull() bool {
	return len(mc.buffer) == mc.maxBufferSize
}

func (mc *MessageChannel) Add(message string) {
	mc.buffer = append(mc.buffer, message)
}

func (mc *MessageChannel) Get() string {
	if len(mc.buffer) == 0 {
		return ""
	}
	message := mc.buffer[0]
	mc.buffer = mc.buffer[1:]
	return message
}

type Producer struct {
	cond           *sync.Cond
	messageChannel *MessageChannel
}

func NewProducer(cond *sync.Cond, messageChannel *MessageChannel) *Producer {
	return &Producer{
		cond:           cond,
		messageChannel: messageChannel,
	}
}

func (p *Producer) Produce(message string) {
	time.Sleep(500 * time.Millisecond) // Simulating some work

	p.cond.L.Lock()
	// Ждем, пока буфер не освободится
	for p.messageChannel.IsFull() {
		fmt.Printf("Producer: жду, буфер полный (%d/%d)\n",
			len(p.messageChannel.buffer), p.messageChannel.maxBufferSize)
		p.cond.Wait()
	}

	p.messageChannel.Add(message)
	fmt.Printf("Producer: добавил '%s', буфер: %d/%d\n",
		message, len(p.messageChannel.buffer), p.messageChannel.maxBufferSize)

	// Сигнализируем потребителю, что появились данные
	p.cond.Signal()
	p.cond.L.Unlock()
}

type Consumer struct {
	cond           *sync.Cond
	messageChannel *MessageChannel
}

func NewConsumer(cond *sync.Cond, messageChannel *MessageChannel) *Consumer {
	return &Consumer{
		cond:           cond,
		messageChannel: messageChannel,
	}
}

func (c *Consumer) Consume() {
	time.Sleep(1 * time.Second) // Simulating some work

	c.cond.L.Lock()
	// Ждем, пока в буфере появятся данные
	for c.messageChannel.IsEmpty() {
		fmt.Println("Consumer: жду, буфер пуст")
		c.cond.Wait()
	}

	message := c.messageChannel.Get()
	fmt.Printf("Consumer: получил '%s', буфер: %d/%d\n",
		message, len(c.messageChannel.buffer), c.messageChannel.maxBufferSize)

	// Сигнализируем производителю, что освободилось место
	c.cond.Signal()
	c.cond.L.Unlock()
}

package main

import (
	"fmt"
	"os"
	"sync"
)

type Messenger struct {
	msg            int
	consumeCond    *sync.Cond
	produceCond    *sync.Cond
	consumersCount int
}

func NewMessenger() *Messenger {
	messenger := Messenger{
		consumeCond: &sync.Cond{
			L: &sync.Mutex{},
		},
		produceCond: &sync.Cond{
			L: &sync.Mutex{},
		},
	}

	return &messenger
}

func (m *Messenger) incConsumersCount() {
	m.produceCond.L.Lock()
	m.consumersCount++
	m.produceCond.Broadcast()
	m.produceCond.L.Unlock()
}

func (m *Messenger) decConsumersCount() {
	m.produceCond.L.Lock()
	m.consumersCount--
	m.produceCond.L.Unlock()
}

func (m *Messenger) waitConsumer() {
	m.produceCond.L.Lock()
	for m.consumersCount == 0 {
		m.produceCond.Wait()
	}
	m.produceCond.L.Unlock()
}

func (m *Messenger) Consume() int {
	m.incConsumersCount()

	m.consumeCond.L.Lock()
	m.consumeCond.Wait()
	msg := m.msg
	m.consumeCond.L.Unlock()

	m.decConsumersCount()

	return msg
}

func (m *Messenger) Produce(msg int) {
	m.waitConsumer()

	m.consumeCond.L.Lock()
	m.msg = msg
	m.consumeCond.Broadcast()
	m.consumeCond.L.Unlock()
}

func main() {
	messenger := NewMessenger()
	var event string
	var msg int
	for {
		fmt.Scanf("%s %d", &event, &msg)
		switch event {
		case "consume":
			fmt.Println("Waiting for message")
			go func() {
				msg := messenger.Consume()
				fmt.Println("Message received: ", msg)
			}()
		case "produce":
			fmt.Println("Generating message")
			go func() {
				messenger.Produce(msg)
			}()
		default:
			fmt.Println("Undefined event")
			os.Exit(1)
		}
	}
}

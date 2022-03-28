package main

import (
	"fmt"
	"sync"
	"time"
)

type (
	Event struct {
		data int
	}

	Observer interface {
		NotifyCallback(Event)
	}

	Subject interface {
		AddListener(Observer)
		RemoveListener(Observer)
		Notify(Event)
	}

	eventObserver struct {
		id   int
		time time.Time
	}

	notificationObserver struct {
		notificationType string
	}

	eventSubject struct {
		observers sync.Map
	}
)

func (e *eventObserver) NotifyCallback(event Event) {
	fmt.Printf("Received: %d after %v\n", event.data, time.Since(e.time))
}

func (e *notificationObserver) NotifyCallback(event Event) {
	fmt.Printf("Sending  %s notification for %d \n", e.notificationType, event.data)
}

func (s *eventSubject) AddListener(obs Observer) {
	fmt.Printf("Adding Listener %+v\n", obs)
	s.observers.Store(obs, struct{}{})
}

func (s *eventSubject) RemoveListener(obs Observer) {
	s.observers.Delete(obs)
}

func (s *eventSubject) Notify(event Event) {
	s.observers.Range(func(key interface{}, value interface{}) bool {
		fmt.Println("Range in sync.Map")
		if key == nil || value == nil {
			return false
		}

		key.(Observer).NotifyCallback(event)
		return true
	})
}

func fib(n int) chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i, j := 0, 1; i < n; i, j = i+j, i {
			out <- i
		}
	}()
	return out
}
func main() {

	n := eventSubject{
		observers: sync.Map{},
	}
	var obs1 = eventObserver{id: 1, time: time.Now()}
	var obs2 = notificationObserver{notificationType: "SMS"}
	var obs3 = notificationObserver{notificationType: "Email"}

	n.AddListener(&obs1)
	n.AddListener(&obs2)
	n.AddListener(&obs3)

	// for x := range fib(3) {
	// 	n.Notify(Event{data: x})
	// 	fmt.Println("\n\n")
	// }

	n.Notify(Event{data: 10})
	fmt.Println()

}

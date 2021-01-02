package tools

import (
	"sync"
	"time"
)

type State struct {
	Notifier sync.WaitGroup
	Finishes sync.WaitGroup
}

func (s *State) WaitUntilAllComplete(timeout time.Duration) (timedOut bool) {
	c := make(chan bool)
	go func() {
		s.Finishes.Wait()
		c <- false
	}()
	go func() {
		time.Sleep(timeout)
		c <- true
	}()
	return <- c
}

func (s *State) AddGoroutine() {
	s.Finishes.Add(1)
}

func (s *State) FinishGoroutine() {
	s.Finishes.Done()
}

func (s *State) WaitUntilShutdownTrigger() {
	s.Notifier.Wait()
}

func (s *State) TriggerShutdown() {
	s.Notifier.Done()
}

func NewState() *State {
	s := new(State)
	s.Notifier.Add(1)
	return s
}
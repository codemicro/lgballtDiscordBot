package state

import (
	"os"
	"os/signal"
	"sync"
	"time"
)

type State struct {
	Notifier       sync.WaitGroup
	Finishes       sync.WaitGroup
	ShutdownSignal chan os.Signal
}

func NewState() *State {
	s := new(State)
	s.ShutdownSignal = make(chan os.Signal, 1)
	signal.Notify(s.ShutdownSignal, os.Interrupt)
	s.Notifier.Add(1)
	return s
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
	return <-c
}

func (s *State) FinishGoroutine() {
	s.Finishes.Done()
}

func (s *State) WaitUntilShutdownTrigger() {
	s.Finishes.Add(1)
	s.Notifier.Wait()
}

func (s *State) TriggerShutdown() {
	s.Notifier.Done()
}

type CustomShutdownSignal struct{}

func (c CustomShutdownSignal) String() string { return "requested shutdown" }
func (c CustomShutdownSignal) Signal()        {}

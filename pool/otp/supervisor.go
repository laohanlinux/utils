package otp

import (
	"errors"
	"sync"
)

// Supervisor as same to erlang supervisor
type Supervisor interface {
	Initization(interface{}) error
	//	NoReply(interface{})
	//	Reply(interface{}) interface{}
	Terminate() error
}

// KGSupervisor waps for kugou util
type KGSupervisor struct {
	onceDo   sync.Once
	currency int
	msg      chan int
}

// Initization with a interface, args is a int slice, args[0] is the currency numbers for children's goroutine
// args[1] is controll children's live type(need or not reload them)
func (s *KGSupervisor) Initization(args interface{}) error {
	if argsValue, ok := args.([]int); ok {
		if len(argsValue) != 2 {
			return errors.New("Invalid Args")
		}
		// start work
		s.onceDo.Do(s.start)
	}
	return nil
}

// NoReply is the channel for user send op command
func (s *KGSupervisor) NoReply(command interface{}) {
	if commandInt, ok := command.(int); ok {
		s.msg <- commandInt
	}
}

// Terminate stop the supervisor children's goroutine
func (s *KGSupervisor) Terminate() {
	s.stop()
}

func (s *KGSupervisor) start() {

}

func (s *KGSupervisor) stop() {

}

// Worker ...
type Worker interface {
	start(...interface{}) error
	terminate()
	reply()
	noReply()
}

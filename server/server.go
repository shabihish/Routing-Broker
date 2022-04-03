package server

import (
	"CA1/helper"
	"math/rand"
	"sync"
	"time"
)

type Msg struct {
	valid bool
	Id    int
	data  string
}

type Server struct {
	currentMessageMu sync.Mutex
	currMessage      Msg
	lastClientId     int
	running          bool
	processing       int
}

func NewMsg(msgId int, data string) *Msg {
	return &Msg{true, msgId, data}
}

func NewServer() *Server {
	return &Server{sync.Mutex{}, Msg{false, 0, ""}, -1, false, 0}
}

type serverInterface interface {
	RunServer()
	getMessages()
	processMessage()
	PutMessage()
	IsRunning()
}

func (s *Server) RunServerAsync() {
	s.running = true

	// Set random num generator seed
	rand.Seed(time.Now().UnixNano())

	// Infinitely look for new messages and process them
	go s.getMessages()
}

func (s *Server) getMessages() {

	idleStatePrinted := false

	for {
		s.currentMessageMu.Lock()
		if !s.currMessage.valid {

			if !idleStatePrinted && s.processing == 0 {
				helper.PrintInColor(helper.ColorBlue, "Server: Now idle!\n")
				idleStatePrinted = true
			}
			s.currentMessageMu.Unlock()
			continue
		}
		idleStatePrinted = false
		if s.processing == 0 {
			helper.PrintInColor(helper.ColorBlue, "Server: Now processing...\n")
		}

		s.processing++
		go s.processMessage()

		s.currMessage.valid = false
		s.currentMessageMu.Unlock()
	}
}

func (s *Server) processMessage() {
	duration := time.Duration((rand.Float32()+0.05)*1000) * time.Millisecond

	helper.PrintInColor(helper.ColorBlue, "Server: Got new message with id %v from client %v to process, will be processed for %v\n", s.currMessage.Id, s.lastClientId, duration)
	time.Sleep(duration)
	s.processing--
}

func (s *Server) PutMessage(msg Msg, clientId int) bool {
	/*
		Returns true if putting the message has been successful, otherwise returns false.
	*/
	s.currentMessageMu.Lock()
	if s.currMessage.valid {
		s.currentMessageMu.Unlock()
		return false
	}

	s.currMessage = msg
	s.lastClientId = clientId
	s.currentMessageMu.Unlock()
	return true
}

func (s *Server) IsRunning() bool {
	return s.running
}

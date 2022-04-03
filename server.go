package main

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

const (
	Mode1Way = 0
	Mode2Way = 1
)

type Msg struct {
	valid bool
	Id    int
	data  string
}

type Server struct {
	currentMessageMu sync.Mutex
	currMessage      Msg
	mode             int
	nextResponseId   int
	lastClientId     int
	running          bool
	processing       int
	brk              *Broker
}

func NewMsg(msgId int, data string) *Msg {
	return &Msg{true, msgId, data}
}

func NewServer() *Server {
	return &Server{sync.Mutex{}, Msg{false, 0, ""}, Mode1Way, 0, -1, false, 0, nil}
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
				PrintLogInColor(ColorBlue, "Server: Now idle!\n")
				idleStatePrinted = true
			}
			s.currentMessageMu.Unlock()
			continue
		}
		idleStatePrinted = false
		if s.processing == 0 {
			PrintLogInColor(ColorBlue, "Server: Now processing...\n")
		}

		s.processing++

		go s.processMessage(s.currMessage.Id, s.lastClientId)
		s.currMessage.valid = false

		s.currentMessageMu.Unlock()
	}
}

func (s *Server) processMessage(msgId int, clientId int) {
	duration := time.Duration((rand.Float32()+0.05)*1000) * time.Millisecond

	PrintLogInColor(ColorBlue, "Server: Got new message with id %v from client %v to process, will be processed for %v\n", msgId, clientId, duration)
	time.Sleep(duration)

	if s.mode == Mode2Way {
		s.generateAndSendResponse(clientId, msgId)
	}
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

func (s *Server) getNextResponseId() int {
	out := s.nextResponseId
	s.nextResponseId++
	return out
}

func (s *Server) generateAndSendResponse(clientId int, responseToMessageId int) {
	msg := *NewMsg(s.getNextResponseId(), "NEXT_RESPONSE")
	s.brk.PutNewServerResponse(msg, clientId, responseToMessageId)
	PrintLogInColor(ColorCyan, "Server: Sent response (%v) for message %v to client %v\n", msg.Id, responseToMessageId, clientId)
}

func (s *Server) Set2WayOn(brk *Broker) error {
	if brk == nil {
		return errors.New("Not a valid Broker instance given!\n")
	}

	s.brk = brk
	s.mode = Mode2Way

	PrintLogInColor(ColorBlue, "Server: Two way response mode turned on!\n")
	return nil
}

func (s *Server) PutAcknowledgement(responseId int, clientId int) {
	PrintLogInColor(ColorCyan, "Server: Received acknowledgment for response (%v) to client %v\n", responseId, clientId)
}


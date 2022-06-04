package server

import (
	"CA1/broker"
	"CA1/helper"
	"errors"
	"math/rand"
	"sync"
	"time"
)

const (
	Mode1Way = 0
	Mode2Way = 1
)

type Server struct {
	currentMessageMu sync.Mutex
	currMessage      helper.Msg
	mode             int
	nextResponseId   int
	lastClientId     int
	running          bool
	processing       int
	brk              *broker.Broker
	logger           *helper.Logger
}

func NewServer(logger *helper.Logger) *Server {
	return &Server{sync.Mutex{}, *helper.NewMsg(false, 0, ""), Mode1Way, 0, -1, false, 0, nil, logger}
}

//type serverInterface interface {
//	RunServer()
//	getMessages()
//	processMessage()
//	PutMessage()
//	IsRunning()
//}

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
		if !s.currMessage.IsValid() {

			if !idleStatePrinted && s.processing == 0 {
				s.logger.PrintLogInColor(helper.ColorBlue, "Server: Now idle!\n")
				idleStatePrinted = true
			}
			s.currentMessageMu.Unlock()
			continue
		}
		idleStatePrinted = false
		if s.processing == 0 {
			s.logger.PrintLogInColor(helper.ColorBlue, "Server: Now processing...\n")
		}

		s.processing++

		go s.processMessage(s.currMessage.Id, s.lastClientId)
		s.currMessage.Invalidate()

		s.currentMessageMu.Unlock()
	}
}

func (s *Server) processMessage(msgId int, clientId int) {
	duration := time.Duration((rand.Float32()+0.05)*1000) * time.Millisecond

	s.logger.PrintLogInColor(helper.ColorBlue, "Server: Got new message with id %v from client %v to process, will be processed for %v\n", msgId, clientId, duration)
	time.Sleep(duration)

	if s.mode == Mode2Way {
		s.generateAndSendResponse(clientId, msgId)
	}
	s.processing--
}

func (s *Server) PutMessage(msg helper.Msg, clientId int) bool {
	/*
		Returns true if putting the message has been successful, otherwise returns false.
	*/
	s.currentMessageMu.Lock()
	if s.currMessage.IsValid() {
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
	msg := *helper.NewMsg(true, s.getNextResponseId(), "NEXT_RESPONSE")
	s.brk.PutNewServerResponse(msg, clientId, responseToMessageId)
	s.logger.PrintLogInColor(helper.ColorCyan, "Server: Sent response (%v) for message %v to client %v\n", msg.Id, responseToMessageId, clientId)
}

func (s *Server) Set2WayOn(brk *broker.Broker) error {
	if brk == nil {
		return errors.New("Not a valid Broker instance given!\n")
	}

	s.brk = brk
	s.mode = Mode2Way

	s.logger.PrintLogInColor(helper.ColorBlue, "Server: Two way response mode turned on!\n")
	return nil
}

func (s *Server) PutAcknowledgement(responseId int, clientId int) {
	s.logger.PrintLogInColor(helper.ColorCyan, "Server: Received acknowledgment for response (%v) to client %v\n", responseId, clientId)
}

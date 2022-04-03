package main

import (
	"errors"
	"math"
	"sync"
)

const (
	ModeSync            = 0
	ModeASync           = 1
	ModeOverflowHandler = 2
	Mode2WayQueue       = 3
)
const DefaultMaxLength = 1

type Broker struct {
	srv        *Server
	mode       int
	queue      queue
	accessList []*Client
}

func NewBroker(srv *Server, mode int) (*Broker, error) {
	if !srv.IsRunning() {
		return nil, errors.New("The server given to this Broker is not running. ")
	}

	maxLength := math.MaxInt64
	if mode == ModeOverflowHandler {
		maxLength = DefaultMaxLength
	}

	is2Way := mode == Mode2WayQueue

	out := &Broker{srv, mode, queue{make([]queueItem, 0), sync.Mutex{}, maxLength, is2Way}, make([]*Client, 0)}
	return out, nil
}

type brokerInterface interface {
	RunBroker()
}

func (brk *Broker) PutNewMessageFromClient(client *Client, msg Msg) {
	if !brk.queue.isFull() {
		if brk.mode == ModeSync && brk.queue.hasClient(client) {
			PrintLogInColor(ColorRed, "Broker: (Sync) Could not push new message (%v) from client %v, client's already got messages to be acknowledged\n", msg.Id, client.ClientId)
			return
		}
		brk.queue.pushBack(queueItem{msg: msg, client: client, isServerResponse: false, responseToMessageId: -1})
		PrintLogInColor(ColorYellow, "Broker: New message (%v) from client %v is added to queue\n", msg.Id, client.ClientId)
	} else {
		PrintLogInColor(ColorRed, "Broker: Could not push new message (%v) from client %v, queue's full\n", msg.Id, client.ClientId)
	}
}

func (brk *Broker) RunBroker() {
	for {
		if brk.queue.isEmpty() {
			continue
		}

		client1, msg, isServerResponse, responseToMessageId := brk.queue.topFront().client, brk.queue.topFront().msg, brk.queue.topFront().isServerResponse, brk.queue.topFront().responseToMessageId
		switch isServerResponse {
		case true:
			// Send response to the corresponding client
			brk.queue.popFront()
			go brk.sendServerResponseToClient(client1, msg, responseToMessageId)
			PrintLogInColor(ColorYellow, "Broker: Sent server response (%v) for message %v to client %v\n", msg.Id, responseToMessageId, client1.ClientId)
		case false:
			// Send client's request to the server
			if brk.srv.PutMessage(msg, client1.ClientId) {
				// Message is now successfully delivered to the server
				PrintLogInColor(ColorYellow, "Broker: Successfully delivered message (%v) from client %v to the server\n", msg.Id, client1.ClientId)

				brk.queue.popFront()
				brk.accessList = append(brk.accessList, client1)
				go client1.PutAcknowledgement(msg.Id)
				PrintLogInColor(ColorYellow, "Broker: Sent acknowledgment for message (%v) to client %v\n", msg.Id, client1.ClientId)
			}
		}
	}

}

func (brk *Broker) findClientById(clientId int) *Client {
	for _, item := range brk.accessList {
		if item.ClientId == clientId {
			return item
		}
	}
	return nil
}

func (brk *Broker) PutNewServerResponse(msg Msg, clientId int, responseToMessageId int) {
	item := queueItem{msg, brk.findClientById(clientId), true, responseToMessageId}
	brk.queue.pushFront(item)
}

func (brk *Broker) sendServerResponseToClient(client *Client, msg Msg, responseToMessageId int) {
	client.PutNewServerResponse(msg, responseToMessageId)
	brk.srv.PutAcknowledgement(msg.Id, client.ClientId)
}

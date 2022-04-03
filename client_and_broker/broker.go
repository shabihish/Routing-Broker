package client_and_broker

import (
	"CA1/helper"
	"CA1/server"
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
const DefaultMaxLength = 8

type Broker struct {
	srv   *server.Server
	mode  int
	queue queue
}

func NewBroker(srv *server.Server, mode int) (*Broker, error) {
	if !srv.IsRunning() {
		return nil, errors.New("The server given to this Broker is not running. ")
	}

	maxLength := math.MaxInt64
	if mode == ModeOverflowHandler {
		maxLength = DefaultMaxLength
	}

	out := &Broker{srv, mode, queue{make([]queueItem, 0), sync.Mutex{}, maxLength, false}}
	return out, nil
}

type brokerInterface interface {
	RunBroker()
}

func (brk *Broker) PutNewMessage(client *Client, msg server.Msg) {
	if !brk.queue.isFull() {
		if brk.mode == ModeSync && brk.queue.hasClient(client){
			helper.PrintInColor(helper.ColorRed, "Broker: (Sync) Could not push new message (%v) from client %v, client's already got messages to be acknowledged\n", msg.Id, client.ClientId)
			return
		}
		brk.queue.push(queueItem{msg: msg, client: client})
		helper.PrintInColor(helper.ColorYellow, "Broker: New message (%v) from client %v is added to queue\n", msg.Id, client.ClientId)
	} else {
		helper.PrintInColor(helper.ColorRed, "Broker: Could not push new message (%v) from client %v, queue's full\n", msg.Id, client.ClientId)
	}
}

func (brk *Broker) RunBroker() {
	for {
		if brk.queue.isEmpty() {
			continue
		}

		client1, msg := brk.queue.top().client, brk.queue.top().msg
		if brk.srv.PutMessage(msg, client1.ClientId) {
			// Message is now successfully delivered to the server
			helper.PrintInColor(helper.ColorYellow, "Broker: Successfully delivered message (%v) from client %v to the server\n", msg.Id, client1.ClientId)

			brk.queue.pop()
			go client1.PutAcknowledgement(msg.Id)
			helper.PrintInColor(helper.ColorYellow, "Broker: Sent acknowledgment for message (%v) to client %v\n", msg.Id, client1.ClientId)
		}
	}

}

package client_and_broker

import (
	"CA1/server"
	"sync"
)

type queueItem struct {
	msg    server.Msg
	client *Client
}

type queue struct {
	queue     []queueItem
	mu        sync.Mutex
	maxLength int
	twoWay    bool
}

type queueInterface interface {
	isEmpty()
	isFull()
	push()
	pop()
	top()
}

func (q *queue) isEmpty() bool {
	q.mu.Lock()
	out := len(q.queue) == 0
	q.mu.Unlock()
	return out
}

func (q *queue) isFull() bool {
	q.mu.Lock()
	out := len(q.queue) >= q.maxLength
	q.mu.Unlock()
	return out
}

func (q *queue) push(item queueItem) bool {
	pushed := false
	q.mu.Lock()
	if len(q.queue) < q.maxLength {
		q.queue = append(q.queue, item)
		pushed = true
	}
	q.mu.Unlock()

	if !pushed {
		return false
	}
	return true
}

func (q *queue) pop() queueItem {
	q.mu.Lock()
	out := q.queue[0]
	q.queue = q.queue[1:]
	q.mu.Unlock()
	return out
}
func (q *queue) top() queueItem {
	q.mu.Lock()
	out := q.queue[0]
	q.mu.Unlock()
	return out
}

func (q *queue) hasClient(client *Client) bool {
	out := false
	q.mu.Lock()
	for _, item := range q.queue {
		if item.client == client {
			out = true
			break
		}
	}
	q.mu.Unlock()
	return out
}

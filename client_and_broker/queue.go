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
	queue []queueItem
	mu sync.Mutex
}

type queueInterface interface {
	isEmpty()
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
func (q *queue) push(item queueItem) {
	q.mu.Lock()
	q.queue = append(q.queue, item)
	q.mu.Unlock()
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

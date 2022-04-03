package main

import (
	"sync"
)

type queueItem struct {
	msg                 Msg
	client              *Client
	isServerResponse    bool
	responseToMessageId int
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

func (q *queue) pushBack(item queueItem) bool {
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

func (q *queue) pushFront(item queueItem) bool {
	pushed := false
	q.mu.Lock()
	if len(q.queue) < q.maxLength {
		q.queue = append([]queueItem{item}, q.queue...)
		pushed = true
	}
	q.mu.Unlock()

	if !pushed {
		return false
	}
	return true
}

func (q *queue) popFront() queueItem {
	q.mu.Lock()
	out := q.queue[0]
	q.queue = q.queue[1:]
	q.mu.Unlock()
	return out
}
func (q *queue) topFront() queueItem {
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

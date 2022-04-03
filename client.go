package main

import (
	"math/rand"
	"time"
)

type Client struct {
	ClientId             int
	acknowledgedMessages map[int]bool
	nextMessageNum       int
	brk                  *Broker
}

func NewClient(clientId int, brk *Broker) *Client {
	if clientId < 0 {
		clientId = -1
	}
	if brk == nil {
		return nil
	}

	return &Client{clientId, make(map[int]bool, 0), 0, brk}
}

type clientInterface interface {
	RunClient()
	generateAndSendNewMessages()
}

func (c *Client) RunClient() {
	for {
		c.generateAndSendNewMessage()
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
}

func (c *Client) getNextMessageId() int {
	out := c.nextMessageNum
	c.nextMessageNum++
	return out
}
func (c *Client) generateAndSendNewMessage() {
	msg := *NewMsg(c.getNextMessageId(), "NEXT_MESSAGE")
	go c.brk.PutNewMessageFromClient(c, msg)

	PrintLogInColor(ColorWhite, "Client %v: Put new message (%v) into the broker queue\n", c.ClientId, msg.Id)
}

func (c *Client) PutAcknowledgement(msgId int) {
	c.acknowledgedMessages[msgId] = true
	PrintLogInColor(ColorGreen, "Client %v: Received acknowledgment for message (%v)\n", c.ClientId, msgId)
}

func (c *Client) PutNewServerResponse(msg Msg, responseToMessageId int) {
	PrintLogInColor(ColorGreen, "Client %v: Received response message (%v) for message %v from the server\n", c.ClientId, msg.Id, responseToMessageId)
}
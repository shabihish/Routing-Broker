package client_and_broker

import (
	"CA1/helper"
	"CA1/server"
	"fmt"
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
		time.Sleep(4 * time.Millisecond)
	}
}

func (c *Client) getNextMessageId() int {
	out := c.nextMessageNum
	c.nextMessageNum++
	return out
}
func (c *Client) generateAndSendNewMessage() {
	msg := *server.NewMsg(c.getNextMessageId(), "NEXT_MESSAGE")
	go c.brk.PutNewMessage(c, msg)

	fmt.Printf("Client %v: Put new message (%v) into the broker queue\n", c.ClientId, msg.Id)
}

func (c *Client) PutAcknowledgement(msgId int) {
	c.acknowledgedMessages[msgId] = true
	helper.PrintInColor(helper.ColorGreen,"Client %v: Received acknowledgment for message (%v)\n", c.ClientId, msgId)
}

package client

import (
	"CA1/broker"
	"CA1/helper"
	"math/rand"
	"time"
)

type Client struct {
	clientId                int
	acknowledgedMessages    map[int]bool
	nextMessageNum          int
	brk                     *broker.Broker
	logger                  *helper.Logger
	assignedClientInterface helper.ClientInterface
}

//type clientInterface interface {
//	RunClient()
//	generateAndSendNewMessages()
//}

func (c *Client) getNextMessageId() int {
	out := c.nextMessageNum
	c.nextMessageNum++
	return out
}

func (c *Client) generateAndSendNewMessage() {
	msg := *helper.NewMsg(true, c.getNextMessageId(), "NEXT_MESSAGE")
	go c.brk.PutNewMessageFromClient(&c.assignedClientInterface, msg)

	c.logger.PrintLogInColor(helper.ColorWhite, "Client %v: Put new message (%v) into the broker queue\n", c.clientId, msg.Id)
}

func NewClient(clientId int, brk *broker.Broker, logger *helper.Logger) *Client {
	if clientId < 0 {
		clientId = -1
	}
	if brk == nil {
		return nil
	}


	c := &Client{clientId, make(map[int]bool, 0), 0, brk, logger, nil}
	c.assignedClientInterface = helper.ClientInterface(c)

	return c
}

func (c *Client) RunClient() {
	for {
		c.generateAndSendNewMessage()
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	}
}

func (c *Client) PutAcknowledgement(msgId int) {
	c.acknowledgedMessages[msgId] = true
	c.logger.PrintLogInColor(helper.ColorGreen, "Client %v: Received acknowledgment for message (%v)\n", c.clientId, msgId)
}

func (c *Client) PutNewServerResponse(msg helper.Msg, responseToMessageId int) {
	c.logger.PrintLogInColor(helper.ColorGreen, "Client %v: Received response message (%v) for message %v from the server\n", c.clientId, msg.Id, responseToMessageId)
}

func (c *Client) GetClientId() int {
	return c.clientId
}

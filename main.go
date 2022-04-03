package main

import (
	"CA1/broker"
	client2 "CA1/client"
	"CA1/helper"
	"CA1/server"
	"sync"
)

const NClients = 10
const Mode = broker.ModeOverflowHandler

var colorMu = sync.Mutex{}

func main() {
	logger := helper.Logger{}

	wg := sync.WaitGroup{}
	wg.Add(NClients + 1)

	srv := server.NewServer(&logger)
	srv.RunServerAsync()

	brk, err := broker.NewBroker(helper.ServerInterface(srv), Mode, &logger)
	if err != nil {
		logger.PrintLogInColor(helper.ColorRed, err.Error()+"\n")
		return
	}

	if Mode == broker.Mode2WayQueue && srv.Set2WayOn(brk) != nil {
		logger.PrintLogInColor(helper.ColorRed, err.Error()+"\n")
		return
	}

	go brk.RunBroker()

	for i := 0; i < NClients; i++ {
		client := client2.NewClient(i, brk, &logger)
		if client == nil {
			return
		}
		go client.RunClient()
	}

	wg.Wait()
}

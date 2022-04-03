package main

import (
	"sync"
)

func printError(err string) {
	PrintLogInColor(ColorRed, "%v\n", err)
}

const NClients = 10
const Mode = Mode2WayQueue

var colorMu = sync.Mutex{}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(NClients + 1)

	srv := NewServer()
	srv.RunServerAsync()

	brk, err := NewBroker(srv, Mode)
	if err != nil {
		printError(err.Error())
		return
	}


	if Mode == Mode2WayQueue && srv.Set2WayOn(brk) != nil {
		printError(err.Error())
		return
	}

	go brk.RunBroker()

	for i := 0; i < NClients; i++ {
		client := NewClient(i, brk)

		if client == nil {
			return
		}
		go client.RunClient()
	}

	wg.Wait()
}

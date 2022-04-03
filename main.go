package main

import (
	"CA1/client_and_broker"
	"CA1/server"
	"fmt"
	"sync"
)

func printError(err string) {
	fmt.Printf("\033[31m%v", err)
}

const NClients = 2

func main() {
	wg := sync.WaitGroup{}
	wg.Add(NClients + 1)

	srv := server.NewServer()
	srv.RunServerAsync()

	brk, err := client_and_broker.NewBroker(srv, client_and_broker.ModeSync)
	if err != nil {
		printError(err.Error())
		return
	}

	go brk.RunBroker()

	for i := 0; i < NClients; i++ {
		client := client_and_broker.NewClient(i, brk)

		if client == nil {
			return
		}
		go client.RunClient()
	}

	wg.Wait()
}

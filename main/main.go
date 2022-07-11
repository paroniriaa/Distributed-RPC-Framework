package main

import (
	"Distributed-RPC-Framework/client"
	"Distributed-RPC-Framework/server"
	"log"
	"net"
	"sync"
	"time"
)

type Demo int

type Input struct {
	Number1 int
	Number2 int
}

func (function Demo) Sum(input Input, output *int) error {
	*output = input.Number1 + input.Number2
	return nil
}

func createServer(address chan string) {
	var demo Demo
	err := server.ServerRegister(&demo)
	if err != nil {
		log.Fatal("Server register error:", err)
	}

	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal("Server Network issue:", err)
	}
	log.Println("RPC server -> createServer: RPC server created and hosting on port", listener.Addr())
	address <- listener.Addr().String()
	server.Connection_handle(listener)
}

func clientCallRPC(client *client.Client, number int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	input := &Input{Number1: number, Number2: number ^ 2}
	var output int
	if err := client.Call("Demo.Sum", input, &output); err != nil {
		log.Fatal("Client RPC call Demo.Sum error:", err)
	}
	log.Printf("%d + %d = %d", input.Number1, input.Number2, output)
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	address := make(chan string)
	go createServer(address)

	testClient, _ := client.MakeDial("tcp", <-address)
	defer func() { _ = testClient.Close() }()

	time.Sleep(time.Second)
	var waitGroup sync.WaitGroup
	n := 0
	for n < 2 {
		waitGroup.Add(1)
		go func(n int) {
			clientCallRPC(testClient, n, &waitGroup)
		}(n)
		n++
	}
	waitGroup.Wait()
}

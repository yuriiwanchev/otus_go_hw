package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Println("Usage: go-telnet --timeout=10s host port")
		return
	}

	address := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1))
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer client.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := client.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "...Connection was closed by peer\n")
			return
		}
	}()

	go func() {
		defer wg.Done()
		if err := client.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "...Connection was closed by peer\n")
			return
		}
	}()

	go func() {
		<-sigChan
		client.Close()
		fmt.Fprintf(os.Stderr, "...Received SIGINT, closing connection\n")
		os.Exit(0)
	}()

	wg.Wait()
	fmt.Fprintf(os.Stderr, "...EOF\n")
}

package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout string
	flag.StringVar(&timeout, "timeout", "10s", "connection timeout")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go-telnet --timeout=10s host port")
		os.Exit(1)
	}

	address := net.JoinHostPort(flag.Args()[0], flag.Args()[1])
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid timeout value")
		os.Exit(1)
	}

	client := NewTelnetClient(address, duration, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		if err := client.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "receive error: %v\n", err)
		}
		os.Exit(0)
	}()

	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "send error: %v\n", err)
		}
		os.Exit(0)
	}()

	select {
	case <-sigCh:
		fmt.Fprintln(os.Stderr, "...Interrupted")
	case <-func() chan struct{} {
		done := make(chan struct{})
		go func() {
			io.Copy(io.Discard, os.Stdin)
			close(done)
		}()
		return done
	}():
		fmt.Fprintln(os.Stderr, "...EOF")
	}
}

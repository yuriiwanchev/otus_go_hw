package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	c.conn = conn
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", c.address)
	return nil
}

func (c *telnetClient) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		return err
	}
	return nil
}

func (c *telnetClient) Send() error {
	if c.conn == nil {
		return fmt.Errorf("connection is closed")
	}
	if c.in == nil {
		return fmt.Errorf("input stream is closed")
	}
	_, err := io.Copy(c.conn, c.in)
	return err
}

func (c *telnetClient) Receive() error {
	if c.conn == nil {
		return fmt.Errorf("connection is closed")
	}
	if c.out == nil {
		return fmt.Errorf("output stream is closed")
	}
	_, err := io.Copy(c.out, c.conn)
	return err
}

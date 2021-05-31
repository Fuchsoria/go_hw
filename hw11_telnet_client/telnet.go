package main

import (
	"bufio"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	c.connection = conn

	return nil
}

func (c *telnetClient) Close() error {
	return c.connection.Close()
}

func (c *telnetClient) Send() error {
	scanner := bufio.NewScanner(c.in)

	for scanner.Scan() {
		_, err := c.connection.Write([]byte(scanner.Text() + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *telnetClient) Receive() error {
	_, err := io.Copy(c.out, c.connection)
	if err != nil {
		return err
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{address, timeout, in, out, nil}
}

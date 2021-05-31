package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout               time.Duration
	ErrWrongArgumentCount = errors.New("wrong arguments count")
	ErrCantConnect        = errors.New("cannot connect to telnet")
)

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout duration")
}

func handleReceiver(errs chan error, client TelnetClient) {
	errs <- client.Receive()
}

func handleSender(errs chan error, client TelnetClient) {
	errs <- client.Send()
}

func main() {
	flag.Parse()

	length := len(os.Args)

	if length != 3 && length != 4 {
		log.Panic(ErrWrongArgumentCount)
	}

	host, port := os.Args[length-2], os.Args[length-1]

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Panic(ErrCantConnect)
	}
	defer client.Close()

	errsCh := make(chan error)
	signCh := make(chan os.Signal, 1)

	signal.Notify(signCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go handleReceiver(errsCh, client)
	go handleSender(errsCh, client)

	select {
	case <-signCh:
		signal.Stop(signCh)
	case <-errsCh:
		return
	}
}

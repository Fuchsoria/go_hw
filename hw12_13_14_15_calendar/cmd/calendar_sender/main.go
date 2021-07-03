package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	simpleconsumer "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/amqp/consumer"
	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar_sender/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	logg := logger.New(config.Logger.Level, config.Logger.File)

	conn, err := amqp.Dial(config.AMPQ.uri)
	if err != nil {
		panic(err)
	}

	c := simpleconsumer.New(config.AMPQ.name, conn, logg)

	msgs, err := c.Consume(ctx, config.AMPQ.name)
	if err != nil {
		logg.Error(fmt.Errorf("cannot consume messages, %w", err).Error())
	}

	logg.Info("start consuming...")

	for m := range msgs {
		fmt.Println("receive new message: ", string(m.Data))
	}

	logg.Info("stopped consuming")
}

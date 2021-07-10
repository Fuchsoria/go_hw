package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	simpleconsumer "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/amqp/consumer"
	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/logger"
	version "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/version"
	"github.com/streadway/amqp"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar_sender/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		version.PrintVersion()
		return
	}

	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logg := logger.New(config.Logger.Level, config.Logger.File)

	conn, err := amqp.Dial(config.AMPQ.URI)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	c := simpleconsumer.New(config.AMPQ.Name, conn, logg)

	msgs, err := c.Consume(ctx, config.AMPQ.Name)
	if err != nil {
		logg.Error(fmt.Errorf("cannot consume messages, %w", err).Error())
	}

	logg.Info("start consuming...")

	for m := range msgs {
		fmt.Println("receive new message: ", string(m.Data))
	}

	logg.Info("stopped consuming")
}

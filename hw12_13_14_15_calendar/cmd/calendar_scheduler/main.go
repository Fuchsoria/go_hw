package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/scheduler"
	sqlstorage "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar_scheduler/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	logg := logger.New(config.Logger.Level, config.Logger.File)

	storage, err := connectStorage(ctx, config)
	if err != nil {
		logg.Error(err.Error())

		panic(err)
	}

	scheduler := scheduler.NewScheduler(logg, storage, config.Scheduler.RecheckDelaySeconds, config.AMPQ.uri, config.AMPQ.name)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		select {
		case <-ctx.Done():
			return
		case <-signals:
		}

		signal.Stop(signals)
		cancel()
	}()

	logg.Info("scheduler is running...")

	scheduler.Start(ctx)
}

func connectStorage(ctx context.Context, config Config) (*sqlstorage.Storage, error) {
	storage, err := sqlstorage.New(ctx, config.DB.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("can't create new storage instance, %w", err)
	}

	err = storage.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't connect to storage, %w", err)
	}

	return storage, nil
}

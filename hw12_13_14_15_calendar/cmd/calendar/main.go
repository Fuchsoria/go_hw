package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/logger"
	gateway "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/server/grpc"
	memorystorage "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage/sql"
	version "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/version"
	_ "github.com/lib/pq"
)

var (
	configFile           string
	ErrCantCreateStorage = errors.New("cannot create storage")
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		version.PrintVersion()

		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logg := logger.New(config.Logger.Level, config.Logger.File)

	storage, err := createStorage(ctx, config)
	if err != nil {
		logg.Error(err.Error())

		log.Fatal(err)
	}

	calendar := app.New(logg, storage)

	server, err := gateway.NewServer(calendar, config.HTTP.Host, config.HTTP.Port, config.HTTP.GrpcPort)
	if err != nil {
		logg.Error(err.Error())
	}

	defer cancel()

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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func createStorage(ctx context.Context, config Config) (app.Storage, error) {
	switch config.DB.Method {
	case "in-memory":
		storage := memorystorage.New()

		return storage, nil
	case "sql":
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

	return nil, ErrCantCreateStorage
}

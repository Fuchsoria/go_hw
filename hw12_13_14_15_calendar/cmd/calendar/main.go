package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/lib/pq"
)

var (
	configFile           string
	ErrCantCreateStorage = errors.New("can not create storage")
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	config := NewConfig()

	logg := logger.New(config.Logger.Level, config.Logger.File)

	storage, err := createStorage(ctx, config)
	if err != nil {
		logg.Error(err.Error())
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(calendar, config.HTTP.Host, config.HTTP.Port)
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
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
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
			return nil, err
		}

		err = storage.Connect(ctx)
		if err != nil {
			return nil, err
		}

		return storage, nil
	}

	return nil, ErrCantCreateStorage
}

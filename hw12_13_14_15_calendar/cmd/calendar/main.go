package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	_ "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/migrations"
	"github.com/pressly/goose"
)

var (
	configFile string
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

	fmt.Println(os.Getwd())
	config := NewConfig()
	fmt.Println(config)
	logg := logger.New(config.Logger.Level, config.Logger.File)

	db, err := goose.OpenDBWithDriver("postgres", config.DB.ConnectionString)
	if err != nil {
		logg.Error(fmt.Sprintf("goose: failed to open DB: %v\n", err))
	}

	defer func() {
		if err := db.Close(); err != nil {
			logg.Error(fmt.Sprintf("goose: failed to close DB: %v\n", err))
		}
	}()

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(calendar)
	ctx, cancel := context.WithCancel(context.Background())
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

	if err := server.Start(ctx, config.HTTP.Host, config.HTTP.Port); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

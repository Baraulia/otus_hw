package main

//nolint:depguard
import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/server/http"
)

var (
	configFile string
	database   string
)

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
	flag.StringVar(&database, "database", "memory", "What database should we use")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.GetLogger(config.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}
	defer logg.Close()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage app.Storage

	switch database {
	case "memory":
		storage = memorystorage.New(logg)
	case "sql":
		storage = sqlstorage.NewPostgresStorage(sqlstorage.PgConfig{
			Host:           config.SQL.Host,
			Username:       config.SQL.Username,
			Password:       config.SQL.Password,
			Port:           config.SQL.Port,
			Database:       config.SQL.Database,
			MigrationsPath: config.SQL.MigrationsPath,
		}, logg)
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: "+err.Error(), nil)
		}

		storage.Close()
	}()

	logg.Info("calendar is running...", nil)

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: "+err.Error(), nil)
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

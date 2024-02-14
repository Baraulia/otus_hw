package main

//nolint:depguard
import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/scheduler"
	sqlstorage "github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/mb"
)

var (
	configFile string
	frequency  int
	database   string
)

func init() {
	flag.StringVar(&configFile, "config", "./configs/scheduler_config.yaml", "Path to configuration file")
	flag.IntVar(&frequency, "frequency", 24, "time between scheduler working(hours)")
	flag.StringVar(&database, "database", "sql", "What database should we use")
}

func main() {
	flag.Parse()

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

	storage := sqlstorage.NewPostgresStorage(sqlstorage.PgConfig{
		Host:           config.SQL.Host,
		Username:       config.SQL.Username,
		Password:       config.SQL.Password,
		Port:           config.SQL.Port,
		Database:       config.SQL.Database,
		MigrationsPath: config.SQL.MigrationsPath,
	}, logg, false)

	connectionURL := fmt.Sprintf("%s://%s:%s@%s:%s",
		config.MB.Protocol, config.MB.Username, config.MB.Password, config.MB.Host, config.MB.Port)

	broker := mb.NewBroker(connectionURL, config.MB.ExchangeName, config.MB.ExchangeType, logg, true)
	err = broker.Connect()
	if err != nil {
		logg.Fatal("error while connecting to message broker", map[string]interface{}{"error": err})
	}

	defer broker.Close()

	if err = broker.InitQueue(config.MB.QueueName, config.MB.RouteKey); err != nil {
		logg.Fatal("error while initialization queue",
			map[string]interface{}{"error": err, "queue name": config.MB.QueueName, "routing key": config.MB.RouteKey})
	}

	producer := mb.NewProducer(broker)

	schedule := scheduler.New(logg, storage, producer, time.Second*time.Duration(frequency), config.MB.RouteKey)
	logg.Info("starting scheduler...", nil)
	go schedule.Start(ctx)

	<-ctx.Done()
	time.Sleep(5 * time.Second)
	logg.Info("closing database...", nil)
	storage.Close()
}

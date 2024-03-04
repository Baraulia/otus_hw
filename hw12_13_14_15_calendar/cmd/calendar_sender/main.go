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

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/sender"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/mb"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender_config.yaml", "Path to configuration file")
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

	connectionURL := fmt.Sprintf("%s://%s:%s@%s:%s",
		config.MB.Protocol, config.MB.Username, config.MB.Password, config.MB.Host, config.MB.Port)

	broker := mb.NewBroker(connectionURL, config.MB.ExchangeName, config.MB.ExchangeType, logg, true)
	err = broker.Connect()
	if err != nil {
		logg.Fatal("error while connecting to message broker", map[string]interface{}{"error": err})
	}
	defer broker.Close()

	if err = broker.InitQueue(config.MB.ConfirmQueueName, config.MB.RouteKey); err != nil {
		logg.Fatal("error while initialization queue",
			map[string]interface{}{"error": err, "queue name": config.MB.ConfirmQueueName, "routing key": config.MB.RouteKey})
	}

	consumer := mb.NewConsumer(config.MB.ClientTag, broker)
	producer := mb.NewProducer(broker)

	notificationSender := sender.New(logg, consumer, producer, config.MB.QueueName, config.MB.RouteKey)
	logg.Info("starting notification sender...", nil)
	go notificationSender.Start(ctx)

	<-ctx.Done()
	logg.Info("closing notification sender...", nil)
	time.Sleep(5 * time.Second)
}

package sender

//nolint:depguard
import (
	"context"
	"fmt"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/mb"
)

type Sender struct {
	logger    Logger
	consumer  mb.ConsumerMB
	routeKey  string
	queueName string
}

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

func New(logger Logger, consumer mb.ConsumerMB, queueName, routeKey string) *Sender {
	return &Sender{logger: logger, consumer: consumer, routeKey: routeKey, queueName: queueName}
}

func (s *Sender) Start(ctx context.Context) {
	s.consumer.ListenQueue(ctx, s.queueName, s.routeKey, func(msg []byte) bool {
		fmt.Println("Notification: ", string(msg))
		return true
	})
}

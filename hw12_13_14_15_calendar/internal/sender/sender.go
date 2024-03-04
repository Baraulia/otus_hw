package sender

//nolint:depguard
import (
	"context"
	"fmt"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/mb"
)

type Sender struct {
	logger           Logger
	producer         mb.ProducerMB
	consumer         mb.ConsumerMB
	routeKey         string
	queueName        string
	confirmQueueName string
}

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

func New(logger Logger, consumer mb.ConsumerMB, producer mb.ProducerMB, queueName, routeKey string) *Sender {
	return &Sender{logger: logger, consumer: consumer, producer: producer, routeKey: routeKey, queueName: queueName}
}

func (s *Sender) Start(ctx context.Context) {
	s.consumer.ListenQueue(ctx, s.queueName, s.routeKey, func(msg []byte) bool {
		fmt.Println("New notification from scheduler: ", string(msg))
		err := s.producer.Publish(s.confirmQueueName, msg)
		if err != nil {
			s.logger.Error("error while sending confirmation", map[string]interface{}{"error": err})
			return false
		}

		return true
	})
}

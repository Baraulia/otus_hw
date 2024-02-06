package mb

//nolint:depguard
import (
	"context"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

type Broker struct {
	connectionURL string
	exchangeName  string
	exchangeType  string
	logger        Logger
	connection    *amqp.Connection
	channel       *amqp.Channel
	done          chan *amqp.Error
	reliable      bool
	confirms      chan amqp.Confirmation
}
type ConsumerMB interface {
	ListenQueue(ctx context.Context, queueName string, routingKey string, handler func([]byte) bool)
}

type ProducerMB interface {
	Publish(routingKey string, msg []byte) error
}

func NewBroker(connectionURL, exchangeName, exchangeType string, logger Logger, reliable bool) *Broker {
	return &Broker{
		connectionURL: connectionURL,
		exchangeName:  exchangeName,
		exchangeType:  exchangeType,
		reliable:      reliable,

		logger: logger,

		done: make(chan *amqp.Error),
	}
}

func (b *Broker) Connect() error {
	var err error
	b.connection, err = amqp.Dial(b.connectionURL)
	if err != nil {
		b.logger.Error("error while dialing to message broker", map[string]interface{}{"error": err})
		return err
	}

	go func() {
		<-b.connection.NotifyClose(b.done)
	}()

	blockings := b.connection.NotifyBlocked(make(chan amqp.Blocking))
	go func() {
		for block := range blockings {
			b.logger.Info(fmt.Sprintf("TCP blocked: %t, reason: %s", block.Active, block.Reason), nil)
		}
	}()

	b.channel, err = b.connection.Channel()
	if err != nil {
		b.logger.Error("error while creating channel", map[string]interface{}{"error": err})
		return err
	}

	err = b.InitExchange(b.exchangeName, b.exchangeType)
	if err != nil {
		return err
	}

	if b.reliable {
		b.logger.Info("enabling publishing confirms", nil)
		if err := b.channel.Confirm(false); err != nil {
			return err
		}

		b.confirms = b.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	}

	return nil
}

func (b *Broker) Close() {
	if b.channel != nil {
		err := b.channel.Close()
		if err != nil {
			b.logger.Error("error while closing channel", map[string]interface{}{"error": err})
			return
		}
		b.channel = nil
	}

	if b.connection != nil {
		err := b.connection.Close()
		if err != nil {
			b.logger.Error("error while closing connection", map[string]interface{}{"error": err})
			return
		}
		b.connection = nil
	}

	if (b.confirms) != nil {
		close(b.confirms)
	}
}

func (b *Broker) Reconnect() error {
	b.Close()
	time.Sleep(5 * time.Second)

	return b.Connect()
}

func (b *Broker) InitExchange(exchangeName, exchangeType string) error {
	err := b.channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		b.logger.Error("error while declaration of exchange", map[string]interface{}{"error": err})
		return err
	}

	return nil
}

func (b *Broker) InitQueue(queueName string, routingKey string) error {
	_, err := b.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		b.logger.Error("error while declaration of queue", map[string]interface{}{"error": err})
		return err
	}

	err = b.channel.QueueBind(queueName, routingKey, b.exchangeName, false, nil)
	if err != nil {
		b.logger.Error("error while queue binding", map[string]interface{}{"error": err})
		return err
	}

	return nil
}

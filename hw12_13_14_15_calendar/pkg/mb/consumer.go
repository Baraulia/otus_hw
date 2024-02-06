package mb

//nolint:depguard
import (
	"context"

	"github.com/streadway/amqp"
)

type Consumer struct {
	clientTag string
	broker    *Broker
}

func NewConsumer(clientTag string, broker *Broker) *Consumer {
	return &Consumer{
		clientTag: clientTag,
		broker:    broker,
	}
}

func (c *Consumer) GetDeliveryChannel(queueName, routingKey string) (<-chan amqp.Delivery, error) {
	if err := c.broker.Connect(); err != nil {
		c.broker.logger.Error("error while connecting to message broker", map[string]interface{}{"error": err})
		return nil, err
	}

	err := c.broker.InitQueue(queueName, routingKey)
	if err != nil {
		c.broker.logger.Error("error while initialization queue",
			map[string]interface{}{"error": err, "queue name": queueName, "routing key": routingKey})
		return nil, err
	}

	deliveryChannel, err := c.broker.channel.Consume(queueName, c.clientTag, false, false, false, false, nil)
	if err != nil {
		c.broker.logger.Error("error while starting delivering queued messages", map[string]interface{}{"error": err})
		return nil, err
	}

	return deliveryChannel, nil
}

func (c *Consumer) ListenQueue(ctx context.Context, queueName string, routingKey string, handler func([]byte) bool) {
	deliveryChannel, err := c.GetDeliveryChannel(queueName, routingKey)
	if err != nil {
		c.broker.logger.Fatal("error while getting delivery channel", map[string]interface{}{"error": err})
		return
	}

	c.Handle(ctx, deliveryChannel, handler, queueName, routingKey)
}

func (c *Consumer) Handle(ctx context.Context, deliveryChannel <-chan amqp.Delivery, handleFunc func([]byte) bool,
	queueName string, routingKey string,
) {
	for {
		select {
		case msg := <-deliveryChannel:
			go c.processMessage(ctx, msg, handleFunc)

		case <-ctx.Done():
			c.broker.logger.Info("closing consumer...", nil)
			c.broker.Close()
			return

		case <-c.broker.done:
			for i := 1; i < 6; i++ {
				c.broker.logger.Info("trying to reconnect to massage broker", map[string]interface{}{"attempt": i})
				err := c.broker.Reconnect()
				if err != nil {
					c.broker.logger.Error("error while reconnecting", map[string]interface{}{"error": err})
				}

				deliveryChannel, err = c.GetDeliveryChannel(queueName, routingKey)
				if err != nil {
					c.broker.logger.Error("error while getting delivery channel", map[string]interface{}{"error": err})
				} else {
					break
				}
			}
			c.broker.logger.Error("reconnect finished", nil)
		}
	}
}

func (c *Consumer) processMessage(_ context.Context, msg amqp.Delivery, handleFunc func([]byte) bool) {
	if ok := handleFunc(msg.Body); ok {
		c.broker.logger.Info("acknowledge message...", nil)
		err := msg.Ack(true)
		if err != nil {
			c.broker.logger.Error("error while acknowledging message",
				map[string]interface{}{"error": err})
		}
	} else {
		err := msg.Nack(false, false)
		if err != nil {
			c.broker.logger.Error("error while negative acknowledging message",
				map[string]interface{}{"error": err})
		}
	}
}

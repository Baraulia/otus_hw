package mb

//nolint:depguard
import (
	"fmt"

	"github.com/streadway/amqp"
)

type Producer struct {
	broker *Broker
}

func NewProducer(broker *Broker) *Producer {
	return &Producer{
		broker,
	}
}

func (p *Producer) Publish(routingKey string, msg []byte) error {
	if p.broker.reliable {
		defer p.confirmOne()
	}

	if err := p.broker.channel.Publish(
		p.broker.exchangeName, // publish to an exchange
		routingKey,            // routing to 0 or more queues
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            msg,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	); err != nil {
		p.broker.logger.Error("error while publishing message", map[string]interface{}{"error": err})
		return err
	}

	return nil
}

func (p *Producer) confirmOne() {
	p.broker.logger.Info("waiting for confirmation of one publishing", nil)
	if confirmed := <-p.broker.confirms; confirmed.Ack {
		p.broker.logger.Info(fmt.Sprintf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag), nil)
	} else {
		p.broker.logger.Info(fmt.Sprintf("failed delivery of delivery tag: %d", confirmed.DeliveryTag), nil)
	}
}

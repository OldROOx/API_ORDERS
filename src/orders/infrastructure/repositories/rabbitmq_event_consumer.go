package repositories

import (
	"app.initial/src/orders/application"
	"app.initial/src/orders/domain/entities"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQEventConsumer struct {
	conn                  *amqp.Connection
	channel               *amqp.Channel
	processPaymentUseCase *application.ProcessPaymentUseCase
}

func NewRabbitMQEventConsumer(
	amqpURL string,
	processPaymentUseCase *application.ProcessPaymentUseCase,
) (*RabbitMQEventConsumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQEventConsumer{
		conn:                  conn,
		channel:               ch,
		processPaymentUseCase: processPaymentUseCase,
	}, nil
}

func (c *RabbitMQEventConsumer) StartConsumingPaymentEvents(queueName, exchange, routingKey string) error {
	// Declare the queue
	q, err := c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// Bind the queue to the exchange
	err = c.channel.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Consume messages
	msgs, err := c.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	go c.handleEvents(msgs)

	return nil
}

func (c *RabbitMQEventConsumer) handleEvents(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		var event entities.Event
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Error deserializing event: %v", err)
			continue
		}

		switch event.Type {
		case "payment.completed":
			orderID, _ := event.Data["order_id"].(float64)
			paymentID, _ := event.Data["payment_id"].(string)

			if err := c.processPaymentUseCase.PaymentCompleted(uint(orderID), paymentID); err != nil {
				log.Printf("Error processing payment completed: %v", err)
			}

		case "payment.failed":
			orderID, _ := event.Data["order_id"].(float64)
			reason, _ := event.Data["reason"].(string)

			if err := c.processPaymentUseCase.PaymentFailed(uint(orderID), reason); err != nil {
				log.Printf("Error processing payment failed: %v", err)
			}
		}
	}
}

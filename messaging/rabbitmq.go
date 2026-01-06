package messaging

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func NewRabbitMQ() (*RabbitMQ, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	exchangeName := "cars_events_exchange"

	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	queueName := "cars_events_queue"

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		"",
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:     conn,
		channel:  ch,
		exchange: exchangeName,
	}, nil
}

func (r *RabbitMQ) PublishEvent(event interface{}) {
	body, err := json.Marshal(event)
	if err != nil {
		log.Println("Ошибка сериализации события:", err)
		return
	}

	err = r.channel.Publish(
		r.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("Ошибка публикации события:", err)
	}
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

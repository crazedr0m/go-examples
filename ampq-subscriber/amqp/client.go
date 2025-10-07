package amqp

import (
	"fmt"

	"ampq_example/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client представляет клиент для работы с RabbitMQ
type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewClient создает нового клиента RabbitMQ
func NewClient(config config.RabbitMQ) (*Client, error) {
	conn, err := amqp.Dial(config.Url)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к AMQP: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Ошибка открытия канала: %s", err)
	}

	// Объявляем exchange
	err = ch.ExchangeDeclare(
		config.ExchangeName, // name
		"direct",            // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("Ошибка объявления exchange: %s", err)
	}

	// Объявляем очередь
	q, err := ch.QueueDeclare(
		config.QueueName, // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("Ошибка объявления очереди: %s", err)
	}

	// Привязываем очередь к exchange
	err = ch.QueueBind(
		q.Name,              // queue name
		config.RoutingKey,   // routing key
		config.ExchangeName, // exchange
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("Ошибка привязки очереди к exchange: %s", err)
	}

	return &Client{
		conn: conn,
		ch:   ch,
	}, nil
}

// Close закрывает соединение с RabbitMQ
func (c *Client) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

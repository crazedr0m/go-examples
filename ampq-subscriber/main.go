// Пример сгенереный нейросетью. теперь надо раздербанить на модули и развиваить дальше
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	yaml "gopkg.in/yaml.v3"
)

type Logger struct {
	level string `yaml:"level"`
}
type RabbitMQ struct {
	url          string `yaml:"url"`
	exchangeName string `yaml:"exchange"`
	queueName    string `yaml:"queue"`
	routingKey   string `yaml:"routing-key"`
	consumerTag  string `yaml:"consumer-tag"`
}

type Config struct {
	log  Logger   `yaml:"log"`
	port int      `yaml:"port"`
	ampq RabbitMQ `yaml:"ampq"`
}

func main() {

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Server Configuration: %+v\n", config)

	// Создаем контекст с отменой для управления жизненным циклом обработчика
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем канал для обработки сигналов операционной системы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем обработчик сообщений в отдельной горутине
	go messageHandler(ctx, config.ampq)

	// Ожидаем сигнал или завершение работы обработчика
	select {
	case <-sigChan:
		log.Println("Получен сигнал, завершение работы...")
		cancel() // Отменяем контекст для завершения работы горутины обработчика
	case <-ctx.Done():
		log.Println("Обработчик завершил работу.")
	}
	log.Println("Программа завершена.")
}

func messageHandler(ctx context.Context, config RabbitMQ) {

	fmt.Println("Обработчик сообщений ")

	conn, err := amqp.Dial(config.url)
	if err != nil {
		log.Fatalf("Ошибка подключения к AMQP: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка открытия канала: %s", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		config.exchangeName, // name
		"direct",            // type
		true,                // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Ошибка объявления exchange: %s", err)
	}

	q, err := ch.QueueDeclare(
		config.queueName, // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Ошибка объявления очереди: %s", err)
	}

	err = ch.QueueBind(
		q.Name,              // queue name
		"",                  // routing key
		config.exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Ошибка привязки очереди к exchange: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,             // queue
		config.consumerTag, // consumer
		false,              // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		log.Fatalf("Ошибка начала потребительской сессии: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			// Обработчик завершает работу при отмене контекста
			return
		case msg := <-msgs:
			log.Printf("Получено сообщение: %s", msg.Body)
			// Обработка сообщения
			// ...
			// Подтверждение получения сообщения (если auto-ack == false)
			err := msg.Ack(false)
			if err != nil {
				log.Printf("Ошибка подтверждения сообщения: %s", err)
			}
		}
	}
}

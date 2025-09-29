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

// Logger представляет конфигурацию логгера
type Logger struct {
	Level string `yaml:"level"`
}

// RabbitMQ представляет конфигурацию подключения к RabbitMQ
type RabbitMQ struct {
	Url          string `yaml:"url"`
	ExchangeName string `yaml:"exchange"`
	QueueName    string `yaml:"queue"`
	RoutingKey   string `yaml:"routing-key"`
	ConsumerTag  string `yaml:"consumer-tag"`
}

// Config представляет общую конфигурацию приложения
type Config struct {
	Log  Logger   `yaml:"log"`
	Port int      `yaml:"port"`
	Ampq RabbitMQ `yaml:"rabbitmq"`
}

// main - точка входа в приложение.
// Читает конфигурацию из файла config.yml и переменных окружения,
// устанавливает обработчик сигналов и запускает обработчик сообщений RabbitMQ.
func main() {

	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Server Configuration: %+v\n", string(data))


	var config Config
	// что-то нихера он не размаршаливает
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

// Читаем переменную окружения
	rabbitUrl := os.Getenv("RABBIT_URL");
	if rabbitUrl != "" {
		config.Ampq.Url = rabbitUrl
	}

	fmt.Printf("Server Configuration: %+v\n", config)

	// Создаем контекст с отменой для управления жизненным циклом обработчика
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем канал для обработки сигналов операционной системы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем обработчик сообщений в отдельной горутине
	go messageHandler(ctx, config.Ampq)

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

// messageHandler обрабатывает сообщения из очереди RabbitMQ.
// Устанавливает соединение с RabbitMQ, объявляет exchange и очередь,
// привязывает очередь к exchange и начинает потреблять сообщения.
// Обработка сообщений продолжается до отмены контекста.
func messageHandler(ctx context.Context, config RabbitMQ) {

	fmt.Println("Обработчик сообщений ")

	conn, err := amqp.Dial(config.Url)
	if err != nil {
		log.Fatalf("Ошибка подключения к AMQP: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка открытия канала: %s", err)
	}
	//	defer ch.Close()

	err = ch.ExchangeDeclare(
		config.ExchangeName, // name
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
		config.QueueName, // name
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
		config.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Ошибка привязки очереди к exchange: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,             // queue
		config.ConsumerTag, // consumer
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

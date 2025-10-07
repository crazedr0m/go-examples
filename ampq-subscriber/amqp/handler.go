package amqp

import (
	"context"
	"fmt"
	"log"
	"time"

	"bytes"
	"os/exec"

	"ampq_example/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageHandler обрабатывает сообщения из очереди RabbitMQ
func MessageHandler(ctx context.Context, client *Client, config config.RabbitMQ) error {
	// Получаем сообщения из очереди
	msgs, err := client.ch.Consume(
		config.QueueName,   // queue
		config.ConsumerTag, // consumer
		false,              // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		return fmt.Errorf("Ошибка начала потребительской сессии: %s", err)
	}

	// Определяем максимальное количество одновременно обрабатываемых сообщений
	maxConcurrency := config.MaxConcurrency
	if maxConcurrency <= 0 {
		maxConcurrency = 5 // Значение по умолчанию
	}

	// Создаем семафор для ограничения количества одновременных обработчиков
	semaphore := make(chan struct{}, maxConcurrency)

	fmt.Printf("Обработчик сообщений запущен с максимальной параллельностью: %d\n", maxConcurrency)

	for {
		select {
		case <-ctx.Done():
			// Обработчик завершает работу при отмене контекста
			return nil
		case msg := <-msgs:
			// Запускаем обработку сообщения в отдельной горутине с ограничением по количеству
			go func(msg amqp.Delivery) {
				// Получаем слот в семафоре
				semaphore <- struct{}{}
				defer func() { <-semaphore }() // Освобождаем слот после завершения

				log.Printf("Получено сообщение: %s", msg.Body)
				log.Printf("длина очереди: %d", len(semaphore))

				//	cmd := exec.Command("php", "-i")
				//	cmd := exec.Command("which", "php")
				//	cmd := exec.Command("ls", "-lah")
				// вот так никогда нельзя делать. огромная дыра в безоасности
				//cmd := exec.Command(string(msg.Body))
				cmd := exec.Command("sh", "./test.sh", string(msg.Body))
				var out bytes.Buffer
				var stderr bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &stderr

				err := cmd.Run()
				if err != nil {
					log.Printf("Command failed: %v\nStderr: %s", err, stderr.String())
					return
				}
				//				fmt.Printf("Output: %s\n", out.String())


				// Обработка сообщения
				time.Sleep(5 * time.Second)
				// ...
				// Подтверждение получения сообщения (если auto-ack == false)
				err = msg.Ack(false)
				if err != nil {
					log.Printf("Ошибка подтверждения сообщения: %s", err)
				}
			}(msg)
		}
	}
}

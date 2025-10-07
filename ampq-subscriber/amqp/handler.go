package amqp

import (
	"context"
	"fmt"
	"log"

	"bytes"
	"os/exec"

	"ampq_example/config"
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

	fmt.Println("Обработчик сообщений запущен")

	for {
		select {
		case <-ctx.Done():
			// Обработчик завершает работу при отмене контекста
			return nil
		case msg := <-msgs:
			log.Printf("Получено сообщение: %s", msg.Body)

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
				log.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
			}
			fmt.Printf("Output: %s\n", out.String())	

			// Обработка сообщения
			// ...
			// Подтверждение получения сообщения (если auto-ack == false)
			err = msg.Ack(false)
			if err != nil {
				log.Printf("Ошибка подтверждения сообщения: %s", err)
			}
		}
	}
}

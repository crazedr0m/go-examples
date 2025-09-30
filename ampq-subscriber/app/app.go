package app

import (
	"context"
	"fmt"
	"log"

	"4v2.com/ampq_example/amqp"
	"4v2.com/ampq_example/config"
)

// App представляет основное приложение
type App struct {
	config *config.Config
	client *amqp.Client
}

// NewApp создает новое приложение
func NewApp(config *config.Config) (*App, error) {
	// Создаем клиент RabbitMQ
	client, err := amqp.NewClient(config.Ampq)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания клиента RabbitMQ: %s", err)
	}

	return &App{
		config: config,
		client: client,
	}, nil
}

// Run запускает приложение
func (a *App) Run(ctx context.Context) error {
	// Запускаем обработчик сообщений в отдельной горутине
	go func() {
		err := amqp.MessageHandler(ctx, a.client, a.config.Ampq)
		if err != nil {
			log.Printf("Ошибка обработчика сообщений: %s", err)
		}
	}()

	// Ожидаем завершения контекста
	<-ctx.Done()
	
	// Закрываем клиент
	a.client.Close()
	
	return nil
}

// Close закрывает приложение
func (a *App) Close() {
	if a.client != nil {
		a.client.Close()
	}
}
package app

import (
	"context"
	"fmt"
	"log"

	"ampq_example/amqp"
	"ampq_example/config"
)

// App представляет основное приложение
type App struct {
	config     *config.Config
	client     *amqp.Client
	configPath string
}

// NewApp создает новое приложение
func NewApp(config *config.Config, configPath string) (*App, error) {
	// Создаем клиент RabbitMQ
	client, err := amqp.NewClient(config.Ampq)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания клиента RabbitMQ: %s", err)
	}

	return &App{
		config:     config,
		client:     client,
		configPath: configPath,
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

// ReloadConfig перезагружает конфигурацию из файла
func (a *App) ReloadConfig() error {
	log.Println("Перезагрузка конфигурации...")
	
	// Загружаем новую конфигурацию
	newConfig, err := config.LoadConfig(a.configPath)
	if err != nil {
		return fmt.Errorf("Ошибка загрузки конфигурации: %s", err)
	}
	
	// Обновляем конфигурацию в приложении
	a.config = newConfig
	log.Println("Конфигурация успешно перезагружена")
	log.Printf("Новая конфигурация: %+v", newConfig)
	
	return nil
}

// Close закрывает приложение
func (a *App) Close() {
	if a.client != nil {
		a.client.Close()
	}
}

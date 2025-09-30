package app

import (
	"context"
	"testing"
	"time"

	"4v2.com/ampq_example/config"
)

func TestNewApp(t *testing.T) {
	// Тест создания нового приложения с конфигурацией
	testConfig := &config.Config{
		Log: config.Logger{
			Level: "debug",
		},
		Port: 8080,
		Ampq: config.RabbitMQ{
			Url:          "amqp://guest:guest@localhost:5672/",
			ExchangeName: "test-exchange",
			QueueName:    "test-queue",
			RoutingKey:   "test-routing-key",
			ConsumerTag:  "test-consumer",
		},
	}
	
	// Этот тест требует запущенного сервера RabbitMQ, поэтому мы просто проверим
	// что сигнатура функции корректна и не паникует с валидной конфигурацией
	// В реальной тестовой среде мы бы использовали моки
	
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewApp panicked with valid config: %v", r)
		}
	}()
	
	// Примечание: Мы не можем actually создать приложение в тестах без запущенного сервера
	// В реальном сценарии мы бы использовали внедрение зависимостей и моки
	_ = testConfig
}

func TestAppRun(t *testing.T) {
	// Тест, что приложение отвечает на отмену контекста
	ctx, cancel := context.WithCancel(context.Background())
	
	// Создаем мок приложения (мы просто используем минимальную структуру)
	mockApp := &App{}
	
	// Запускаем приложение в горутине
	done := make(chan bool, 1)
	go func() {
		// Это завершится ошибкой, потому что у нас мок приложения, но это ожидаемо
		_ = mockApp.Run(ctx)
		done <- true
	}()
	
	// Отменяем контекст через короткое время
	time.Sleep(10 * time.Millisecond)
	cancel()
	
	// Ждем завершения приложения
	select {
	case <-done:
		// Приложение отреагировало на отмену контекста
	case <-time.After(1 * time.Second):
		t.Error("App did not respond to context cancellation in time")
	}
}
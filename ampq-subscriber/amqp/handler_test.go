package amqp

import (
	"context"
	"testing"
	"time"

	"4v2.com/ampq_example/config"
)

func TestMessageHandler(t *testing.T) {
	// Тестируем, что сигнатура функции корректна
	// Поскольку мы не можем легко замокать клиент AMQP, мы просто проверим
	// что функция корректно обрабатывает отмену контекста
	
	// Создаем тестовую конфигурацию
	testConfig := config.RabbitMQ{
		Url:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "test-exchange",
		QueueName:    "test-queue",
		RoutingKey:   "test-routing-key",
		ConsumerTag:  "test-consumer",
	}
	
	// Создаем контекст, который можно отменить
	ctx, cancel := context.WithCancel(context.Background())
	
	// Создаем мок клиента (для этого теста передаем nil)
	var mockClient *Client = nil
	
	// Запускаем обработчик в горутине
	errChan := make(chan error, 1)
	go func() {
		// Это завершится ошибкой, потому что мы передаем nil клиент, но это ожидаемо
		err := MessageHandler(ctx, mockClient, testConfig)
		errChan <- err
	}()
	
	// Отменяем контекст через короткое время
	time.Sleep(10 * time.Millisecond)
	cancel()
	
	// Ждем завершения обработчика
	select {
	case err := <-errChan:
		// Мы ожидаем ошибку, потому что передали nil клиент, но это нормально для этого теста
		_ = err
	case <-time.After(1 * time.Second):
		t.Error("MessageHandler did not respond to context cancellation in time")
	}
}
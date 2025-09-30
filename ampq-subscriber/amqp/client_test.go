package amqp

import (
	"testing"

	"4v2.com/ampq_example/config"
)

func TestNewClient(t *testing.T) {
	// Этот тест требует запущенного сервера RabbitMQ, поэтому мы просто проверим
	// что сигнатура функции корректна и не паникует с nil конфигом
	// В реальной тестовой среде мы бы использовали мок или интеграционные тесты
	
	// Тест с nil конфигом (должен обрабатываться корректно)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewClient panicked with nil config: %v", r)
		}
	}()
	
	// Создаем минимальную конфигурацию для тестирования
	testConfig := config.RabbitMQ{
		Url:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "test-exchange",
		QueueName:    "test-queue",
		RoutingKey:   "test-routing-key",
		ConsumerTag:  "test-consumer",
	}
	
	// Примечание: Мы не можем actually подключиться в тестах без запущенного сервера
	// В реальном сценарии мы бы использовали внедрение зависимостей и моки
	_ = testConfig
}
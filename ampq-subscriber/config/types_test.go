package config

import (
	"testing"
)

func TestConfigStructs(t *testing.T) {
	// Тест структуры Logger
	logger := Logger{
		Level: "info",
	}
	if logger.Level != "info" {
		t.Errorf("Expected Level to be 'info', got %s", logger.Level)
	}

	// Тест структуры RabbitMQ
	rabbitMQ := RabbitMQ{
		Url:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "test-exchange",
		QueueName:    "test-queue",
		RoutingKey:   "test-routing-key",
		ConsumerTag:  "test-consumer",
	}
	if rabbitMQ.Url != "amqp://guest:guest@localhost:5672/" {
		t.Errorf("Expected Url to be 'amqp://guest:guest@localhost:5672/', got %s", rabbitMQ.Url)
	}

	// Тест структуры Config
	config := Config{
		Log: Logger{
			Level: "debug",
		},
		Port: 8080,
		Ampq: RabbitMQ{
			Url:          "amqp://guest:guest@localhost:5672/",
			ExchangeName: "test-exchange",
			QueueName:    "test-queue",
			RoutingKey:   "test-routing-key",
			ConsumerTag:  "test-consumer",
		},
	}
	if config.Log.Level != "debug" {
		t.Errorf("Expected Log.Level to be 'debug', got %s", config.Log.Level)
	}
	if config.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", config.Port)
	}
}

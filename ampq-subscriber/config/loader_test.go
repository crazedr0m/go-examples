package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Создаем временный файл конфигурации для тестирования
	tempConfig := `log:
  level: debug
port: 8080
rabbitmq:
  url: amqp://guest:guest@localhost:5672/
  exchange: test-exchange
  queue: test-queue
  routing-key: test-routing-key
  consumer-tag: test-consumer
`
	
	// Записываем временный файл конфигурации
	err := os.WriteFile("test_config.yml", []byte(tempConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}
	
	// Очищаем после теста
	defer os.Remove("test_config.yml")
	
	// Тестируем загрузку конфигурации
	config, err := LoadConfig("test_config.yml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Тестируем значения конфигурации
	if config.Log.Level != "debug" {
		t.Errorf("Expected Log.Level to be 'debug', got %s", config.Log.Level)
	}
	
	if config.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", config.Port)
	}
	
	if config.Ampq.Url != "amqp://guest:guest@localhost:5672/" {
		t.Errorf("Expected Ampq.Url to be 'amqp://guest:guest@localhost:5672/', got %s", config.Ampq.Url)
	}
	
	if config.Ampq.ExchangeName != "test-exchange" {
		t.Errorf("Expected Ampq.ExchangeName to be 'test-exchange', got %s", config.Ampq.ExchangeName)
	}
	
	if config.Ampq.QueueName != "test-queue" {
		t.Errorf("Expected Ampq.QueueName to be 'test-queue', got %s", config.Ampq.QueueName)
	}
}

func TestLoadConfigWithEnv(t *testing.T) {
	// Создаем временный файл конфигурации для тестирования
	tempConfig := `log:
  level: debug
port: 8080
rabbitmq:
  url: amqp://guest:guest@localhost:5672/
  exchange: test-exchange
  queue: test-queue
  routing-key: test-routing-key
  consumer-tag: test-consumer
`
	
	// Записываем временный файл конфигурации
	err := os.WriteFile("test_config_env.yml", []byte(tempConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}
	
	// Очищаем после теста
	defer os.Remove("test_config_env.yml")
	
	// Устанавливаем переменную окружения
	os.Setenv("RABBIT_URL", "amqp://test:test@localhost:5672/")
	defer os.Unsetenv("RABBIT_URL")
	
	// Тестируем загрузку конфигурации
	config, err := LoadConfig("test_config_env.yml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Тестируем, что переменная окружения переопределяет конфигурацию
	if config.Ampq.Url != "amqp://test:test@localhost:5672/" {
		t.Errorf("Expected Ampq.Url to be 'amqp://test:test@localhost:5672/' from env, got %s", config.Ampq.Url)
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	// Тестируем загрузку несуществующего файла конфигурации
	_, err := LoadConfig("non_existent_config.yml")
	if err == nil {
		t.Error("Expected error when loading non-existent config file, but got nil")
	}
}
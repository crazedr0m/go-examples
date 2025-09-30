package config

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
Вот пример Golang-программы, которая подключается к AMQP-брокеру, использует настройки из файла конфигурации, асинхронно обрабатывает сообщения и работает как демон:
1. Настройка проекта
Создайте каталог для вашего проекта (например, amqp-consumer) и перейдите в него.
Инициализируйте Go-модуль: go mod init amqp-consumer.
Установите необходимые библиотеки:
go get github.com/rabbitmq/amqp091-go для работы с RabbitMQ.
go get github.com/spf13/viper для работы с конфигурационными файлами.
go get github.com/sirupsen/logrus для логирования.
2. Файл конфигурации (config.yaml)
Создайте файл config.yaml в корне вашего проекта со следующим содержимым:
amqp:
  url: amqp://guest:guest@localhost:5672/
  exchange: my-exchange
  queue: my-queue
  routing-key: my-routing-key
  log:
  level: info
3. Программа (main.go)
Создайте файл main.go в корне вашего проекта со следующим содержимым:
package main

import (
    "context"
    "log"

    "github.com/rabbitmq/amqp091-go"
    "github.com/spf13/viper"
)

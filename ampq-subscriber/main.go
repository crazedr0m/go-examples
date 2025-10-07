// Пример сгенереный нейросетью. теперь надо раздербанить на модули и развиваить дальше
package main

import (
	"context"
	"fmt"
	"log"

	"ampq_example/app"
	"ampq_example/config"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	log.Printf("Server Configuration: %+v\n", cfg)

	// Создаем контекст с отменой для управления жизненным циклом обработчика
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настраиваем обработчик сигналов
	// Создаем приложение
	application, err := app.NewApp(cfg, "config.yml")
	if err != nil {
		log.Fatalf("Ошибка создания приложения: %v", err)
	}
	defer application.Close()

	// Настраиваем обработчик сигналов
	app.SetupSignalHandler(ctx, cancel, application)

	fmt.Printf("Server Configuration: %+v\n", cfg)

	// Запускаем приложение
	err = application.Run(ctx)
	if err != nil {
		log.Fatalf("Ошибка запуска приложения: %v", err)
	}

	log.Println("Программа завершена.")
}

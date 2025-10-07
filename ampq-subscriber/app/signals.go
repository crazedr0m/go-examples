package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// SetupSignalHandler настраивает обработчик сигналов операционной системы
func SetupSignalHandler(ctx context.Context, cancel context.CancelFunc, app *App) {
	// Создаем канал для обработки сигналов операционной системы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sigNotice := make(chan os.Signal, 1)
	signal.Notify(sigNotice, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)

	// Запускаем горутину для обработки сигналов
	go func() {
		for {
			select {
			case sig := <-sigNotice:
				switch sig {
				case syscall.SIGUSR1:
					if app != nil {
						log.Println("Получен сигнал SIGUSR1, перезагрузка конфигурации...")
						if err := app.ReloadConfig(); err != nil {
							log.Printf("Ошибка перезагрузки конфигурации: %s", err)
						}
					} else {
						log.Println("Получен сигнал SIGUSR1, но App не инициализирован")
					}
				default:
					log.Printf("Получен сигнал %s, но обработка не реализована", sig)
				}
			case <-sigChan:
				log.Println("Получен сигнал, завершение работы...")
				cancel() // Отменяем контекст для завершения работы
				return
			case <-ctx.Done():
				// Контекст уже отменен, выходим из горутины
				return
			}
		}
	}()
}
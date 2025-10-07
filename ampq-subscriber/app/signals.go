package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// SetupSignalHandler настраивает обработчик сигналов операционной системы
func SetupSignalHandler(ctx context.Context, cancel context.CancelFunc) {
	// Создаем канал для обработки сигналов операционной системы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sigNotice := make(chan os.Signal, 1)
	signal.Notify(sigNotice, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)

	// Запускаем горутину для обработки сигналов
	go func() {
		for {
			select {
			case signal := <-sigNotice:
				log.Printf("Получен сигнал %s, что-то надо сделать...\n", signal)
			case signal := <-sigChan:
				log.Printf("Получен сигнал %s, завершение работы...\n", signal)
				cancel() // Отменяем контекст для завершения работы
				return
			case <-ctx.Done():
			// Контекст уже отменен, ничего не делаем
				return
			}
		}
	}()
}
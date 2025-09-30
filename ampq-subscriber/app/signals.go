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

	// Запускаем горутину для обработки сигналов
	go func() {
		select {
		case <-sigChan:
			log.Println("Получен сигнал, завершение работы...")
			cancel() // Отменяем контекст для завершения работы
		case <-ctx.Done():
			// Контекст уже отменен, ничего не делаем
		}
	}()
}
package app

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestSetupSignalHandler(t *testing.T) {
	// Тест, что обработчик сигналов отвечает на SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	
	// Настраиваем обработчик сигналов
	SetupSignalHandler(ctx, cancel)
	
	// Отправляем сигнал SIGTERM процессу
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find current process: %v", err)
	}
	
	// Отправляем сигнал в горутине
	signalSent := make(chan bool, 1)
	go func() {
		err := proc.Signal(syscall.SIGTERM)
		if err != nil {
			t.Errorf("Failed to send signal: %v", err)
		}
		signalSent <- true
	}()
	
	// Ждем либо отмены контекста, либо таймаута
	select {
	case <-ctx.Done():
		// Контекст был отменен, что мы и ожидали
		<-signalSent // Wait for signal to be sent
	case <-time.After(2 * time.Second):
		t.Error("Signal handler did not cancel context in time")
	}
}
package app

import (
	"context"
	"flag"
	"log/slog"
	"reflect"
	"testing"
)

func TestSetupLogger(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		wantHandler string
		wantLevel   slog.Level
	}{
		{
			name:        "envLocal",
			envValue:    "local",
			wantHandler: "*slog.TextHandler", // Для 'local' ожидаем TextHandler
			wantLevel:   slog.LevelDebug,     // Уровень debug для локального окружения
		},
		{
			name:        "envDev",
			envValue:    "dev",
			wantHandler: "*slog.JSONHandler", // Для 'dev' ожидаем JSONHandler
			wantLevel:   slog.LevelDebug,     // Уровень debug для dev
		},
		{
			name:        "envProd",
			envValue:    "prod",
			wantHandler: "*slog.JSONHandler", // Для 'prod' JSONHandler
			wantLevel:   slog.LevelInfo,      // Уровень info для prod
		},
		{
			name:        "unknown env",
			envValue:    "unknown",
			wantHandler: "*slog.JSONHandler", // Для неизвестного окружения JSONHandler
			wantLevel:   slog.LevelInfo,      // Уровень info по умолчанию
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем флаг
			if err := flag.Set("env", tt.envValue); err != nil {
				t.Error("set flag error:", err)
			}
			// Сбрасываем парсинг флагов
			flag.Parse()

			// Вызываем SetupLogger
			got := SetupLogger()

			// Проверяем тип обработчика
			handlerType := reflect.TypeOf(got.Handler()).String()
			if handlerType != tt.wantHandler {
				t.Errorf("SetupLogger() handler = %v, want %v", handlerType, tt.wantHandler)
			}

			// Проверка уровня логирования. Здесь мы не можем напрямую получить уровень,
			// но можно проверить, что нужный уровень включен через Enabled().
			if tt.wantLevel == slog.LevelDebug && !got.Handler().Enabled(context.Background(), slog.LevelDebug) {
				t.Errorf("Expected Debug level to be enabled for %v", tt.wantHandler)
			}

			if tt.wantLevel == slog.LevelInfo && !got.Handler().Enabled(context.Background(), slog.LevelInfo) {
				t.Errorf("Expected Info level to be enabled for %v", tt.wantHandler)
			}
		})
	}
}

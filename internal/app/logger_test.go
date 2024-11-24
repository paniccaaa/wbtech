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
			wantHandler: "*slog.TextHandler",
			wantLevel:   slog.LevelDebug,
		},
		{
			name:        "envDev",
			envValue:    "dev",
			wantHandler: "*slog.JSONHandler",
			wantLevel:   slog.LevelDebug,
		},
		{
			name:        "envProd",
			envValue:    "prod",
			wantHandler: "*slog.JSONHandler",
			wantLevel:   slog.LevelInfo,
		},
		{
			name:        "unknown env",
			envValue:    "unknown",
			wantHandler: "*slog.JSONHandler",
			wantLevel:   slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := flag.Set("env", tt.envValue); err != nil {
				t.Error("set flag error:", err)
			}
			flag.Parse()

			got := SetupLogger()

			handlerType := reflect.TypeOf(got.Handler()).String()
			if handlerType != tt.wantHandler {
				t.Errorf("SetupLogger() handler = %v, want %v", handlerType, tt.wantHandler)
			}

			if tt.wantLevel == slog.LevelDebug && !got.Handler().Enabled(context.Background(), slog.LevelDebug) {
				t.Errorf("Expected Debug level to be enabled for %v", tt.wantHandler)
			}

			if tt.wantLevel == slog.LevelInfo && !got.Handler().Enabled(context.Background(), slog.LevelInfo) {
				t.Errorf("Expected Info level to be enabled for %v", tt.wantHandler)
			}
		})
	}
}

package dev01

import (
	"log/slog"
	"os"
	"time"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func Now() time.Time {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		slog.Error("failed to get time from ntp server", slog.String("err", err.Error()))
		os.Exit(1)
	}

	return t
}

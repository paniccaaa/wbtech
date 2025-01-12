package dev01

import (
	"fmt"
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

type S struct {
	t time.Time
}

func Now() (S, error) {
	var res S

	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return S{}, err
	}

	res.t = t

	fmt.Println("ntp", t)
	fmt.Println(time.Now())

	return res, nil
}

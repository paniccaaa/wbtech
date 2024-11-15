package model

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrOrderNotFound = fmt.Errorf("order not found: %w", ErrNotFound)
)

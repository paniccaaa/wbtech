package postgres

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/paniccaaa/wbtech/internal/model"
	"github.com/paniccaaa/wbtech/internal/services/order"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(DB_URI string) order.Storage {
	db, err := sqlx.Connect("postgres", DB_URI)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to verify connection to db: %v", err)
	}

	return &Repository{db: db}
}

func (r *Repository) GetOrder(ctx context.Context, order_uid model.OrderUID) (model.Order, error) {
	return model.Order{}, nil
}

func (r *Repository) SaveOrder(ctx context.Context, order model.Order) error {
	return nil
}

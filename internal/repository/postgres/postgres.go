package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/paniccaaa/wbtech/internal/model"
	"github.com/paniccaaa/wbtech/internal/services/order"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(DB_URI string) (order.Storage, error) {
	db, err := sqlx.Connect("postgres", DB_URI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to verify connection to db: %v", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) GetOrder(ctx context.Context, order_uid model.OrderUID) (model.Order, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Select("order_data").
		From("orders").
		Where(sq.And{
			sq.Eq{"order_uid": string(order_uid)},
		})

	q, args, err := query.ToSql()
	if err != nil {
		return model.Order{}, fmt.Errorf("build sql: %w", err)
	}

	var orderData string
	if err := r.db.GetContext(ctx, &orderData, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}

		return model.Order{}, fmt.Errorf("get order: %w", err)
	}

	var order model.Order
	if err := json.Unmarshal([]byte(orderData), &order); err != nil {
		return model.Order{}, fmt.Errorf("unmarshal order: %w", err)
	}

	return order, nil
}

func (r *Repository) SaveOrder(ctx context.Context, order model.Order) error {
	orderData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Insert("orders").
		Columns("order_uid", "order_data").
		Values(order.OrderUID, orderData).
		Suffix("ON CONFLICT (order_uid) DO UPDATE SET order_data = EXCLUDED.order_data")

	q, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build sql: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("save order: %w", err)
	}

	return nil
}

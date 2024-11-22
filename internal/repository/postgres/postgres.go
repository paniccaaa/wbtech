package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/paniccaaa/wbtech/internal/model"
	"github.com/paniccaaa/wbtech/internal/services/order"
)

type Cache struct {
	mu    sync.RWMutex
	cache map[model.OrderUID]model.Order
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[model.OrderUID]model.Order),
	}
}

func (c *Cache) Get(orderUID model.OrderUID) (model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.cache[orderUID]
	return order, ok
}

func (c *Cache) Set(order model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[order.OrderUID] = order
}

func (c *Cache) Restore(orders []model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		c.cache[order.OrderUID] = order
	}

}

type Repository struct {
	db    *sqlx.DB
	cache *Cache
	log   *slog.Logger
}

func NewRepository(DB_URI string, log *slog.Logger) (order.Storage, error) {
	db, err := sqlx.Connect("postgres", DB_URI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}

	log.Info("connected to db", slog.String("DB_URI", DB_URI))

	cache := NewCache()
	repo := &Repository{db: db, cache: cache, log: log}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to verify connection to db: %v", err)
	}

	if err := repo.restoreCacheFromDB(); err != nil {
		return nil, fmt.Errorf("restore cache: %w", err)
	}

	log.Info("successfully restore cache")

	return repo, nil
}

func (r *Repository) restoreCacheFromDB() error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("order_data").From("orders")

	q, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build sql: %w", err)
	}

	var ordersData []string
	err = r.db.Select(&ordersData, q, args...)
	if err != nil {
		return fmt.Errorf("restore cache: %w", err)
	}

	var orders []model.Order
	for _, orderJSON := range ordersData {
		var order model.Order
		if err := json.Unmarshal([]byte(orderJSON), &order); err != nil {
			return fmt.Errorf("unmarshal order: %w", err)
		}
		orders = append(orders, order)
	}

	r.cache.Restore(orders)
	return nil
}

func (r *Repository) GetOrder(ctx context.Context, orderUID model.OrderUID) (model.Order, error) {
	if order, found := r.cache.Get(orderUID); found {
		r.log.Info("get from cache", slog.String("orderUID", string(orderUID)))
		return order, nil
	}

	r.log.Info("order not found", slog.String("orderUID", string(orderUID)))

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Select("order_data").
		From("orders").
		Where(sq.Eq{"order_uid": string(orderUID)})

	q, args, err := query.ToSql()
	if err != nil {
		return model.Order{}, fmt.Errorf("build sql: %w", err)
	}

	var orderData string
	if err := r.db.GetContext(ctx, &orderData, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}

		r.log.Error("failed to get order", slog.String("err", err.Error()))

		return model.Order{}, fmt.Errorf("get order: %w", err)
	}

	var order model.Order
	if err := json.Unmarshal([]byte(orderData), &order); err != nil {
		return model.Order{}, fmt.Errorf("unmarshal order: %w", err)
	}

	r.cache.Set(order)
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

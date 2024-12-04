package solutions

import (
	"fmt"
	"log/slog"
)

// DONE
// Реализовать паттерн «адаптер» на любом примере.

type Storage interface {
	Insert(order []byte) error
	GetOrderByID(id string) ([]byte, error)
}

func service(db Storage) {
	_ = db.Insert([]byte("some order"))

	order, _ := db.GetOrderByID("1")
	fmt.Println("service:", order)
}

type RepositoryPG struct {
}

func (r *RepositoryPG) Insert(order []byte) error {
	slog.Info("successfully save order in postgres")
	return nil
}

func (r *RepositoryPG) GetOrderByID(id string) ([]byte, error) {
	slog.Info("successfully get order by id")
	return []byte("order data"), nil
}

type RepositoryMongo struct {
}

func (r *RepositoryMongo) Save(order string) error {
	slog.Info("successfully save order in mongo")
	return nil
}

func (r *RepositoryMongo) FindOrder(id string) (string, error) {
	slog.Info("successfully find order in mongo")
	return "order data from mongo", nil
}

type MongoAdapter struct {
	mongo *RepositoryMongo
}

func (m *MongoAdapter) Insert(order []byte) error {
	return m.mongo.Save(string(order))
}

func (m *MongoAdapter) GetOrderByID(id string) ([]byte, error) {
	order, err := m.mongo.FindOrder(id)
	if err != nil {
		return nil, fmt.Errorf("find order: %w", err)
	}

	return []byte(order), err
}

func Solve21() {
	// usage pg
	postgres := &RepositoryPG{}
	service(postgres)

	// usage mongo
	mongo := &RepositoryMongo{}
	mongoAdapter := &MongoAdapter{mongo: mongo}
	service(mongoAdapter)
}

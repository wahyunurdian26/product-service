package postgres

import (
	"database/sql"

	"github.com/wahyunurdian26/product-service/repository"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

package postgres

import (
	"context"

	"github.com/wahyunurdian26/product-service/model"
)

var (
	createProductQuery = `
		INSERT INTO products (id, name, price, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
)

func (r *productRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	_, err := r.db.ExecContext(ctx, createProductQuery,
		product.ID,
		product.Name,
		product.Price,
		product.Type,
		product.CreatedAt,
		product.UpdatedAt,
	)
	return err
}

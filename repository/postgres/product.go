package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/wahyunurdian26/product-service/model"
	"github.com/wahyunurdian26/product-service/repository"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (id, name, price, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		product.ID,
		product.Name,
		product.Price,
		product.Type,
		product.CreatedAt,
		product.UpdatedAt,
	)
	return err
}

func (r *productRepository) ListProducts(ctx context.Context, req model.ListProductRequest) ([]model.Product, error) {
	query := `SELECT id, name, price, type, created_at, updated_at FROM products WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if req.Search != "" {
		query += fmt.Sprintf(` AND (name ILIKE $%d OR id = $%d)`, argIdx, argIdx+1)
		args = append(args, "%"+req.Search+"%", req.Search)
		argIdx += 2
	}

	if req.Type != "" {
		query += fmt.Sprintf(` AND type = $%d`, argIdx)
		args = append(args, req.Type)
		argIdx++
	}

	// Sorting
	sortBy := "created_at"
	if req.SortBy == "price" || req.SortBy == "name" || req.SortBy == "created_at" {
		sortBy = req.SortBy
	}

	order := "DESC"
	if strings.ToUpper(req.Order) == "ASC" {
		order = "ASC"
	}

	query += fmt.Sprintf(` ORDER BY %s %s`, sortBy, order)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Type, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/wahyunurdian26/product-service/model"
)

var (
	baseListProductsQuery = `SELECT id, name, price, type, created_at, updated_at FROM products WHERE 1=1`
)

func (r *productRepository) ListProducts(ctx context.Context, req model.ListProductRequest) ([]model.Product, error) {
	query := baseListProductsQuery
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

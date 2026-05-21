package model

import (
	"time"
)

type ProductType string

const (
	TypeSayuran ProductType = "Sayuran"
	TypeProtein ProductType = "Protein"
	TypeBuah    ProductType = "Buah"
	TypeSnack   ProductType = "Snack"
)

type Product struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Price     float64     `json:"price"`
	Type      ProductType `json:"type"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CreateProductRequest struct {
	Name  string      `json:"name"`
	Price float64     `json:"price"`
	Type  ProductType `json:"type"`
}

type CreateProductResponse struct {
	Product *Product `json:"product,omitempty"`
	Message string   `json:"message,omitempty"`
	Error   string   `json:"error,omitempty"`
}

type ListProductRequest struct {
	Search string `json:"search"`
	Type   string `json:"type"`
	SortBy string `json:"sort_by"`
	Order  string `json:"order"`
}

type ListProductResponse struct {
	Products []Product `json:"products"`
	Message  string    `json:"message,omitempty"`
	Error    string    `json:"error,omitempty"`
}

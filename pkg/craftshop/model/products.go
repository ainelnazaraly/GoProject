package model

/// we define API’s model
import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ainelnazaraly/CraftShop/pkg/craftshop/validator"
)

type Product struct {
	ID              int     `json:"product_id"`
	SellerID        int     `json:"seller_id"`
	Name            string  `json:"product_name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Category        string  `json:"category"`
	MaterialsUsed   string  `json:"materials_used"`
	ShippingDetails string  `json:"shipping_details"`
}

type ProductModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func ValidateProduct(v *validator.Validator, product *Product) {
	// Check if the product name field is empty.
	v.Check(product.Name != "", "product_name", "must be provided")
	// Check if the product name field is not more than 100 characters.
	v.Check(len(product.Name) <= 100, "product_name", "must not be more than 100 characters long")
	// Check if the description field is empty.
	v.Check(product.Description != "", "description", "must be provided")
	// Check if the price field is greater than zero.
	v.Check(product.Price > 0, "price", "must be greater than zero")
	// Check if the category field is empty.
	v.Check(product.Category != "", "category", "must be provided")
	// Check if the category field is not more than 100 characters.
	v.Check(len(product.Category) <= 100, "category", "must not be more than 100 characters long")
	// Check if the materials used field is not more than 255 characters.
	v.Check(len(product.MaterialsUsed) <= 255, "materials_used", "must not be more than 255 characters long")
	// Check if the shipping details field is not more than 255 characters.
	v.Check(len(product.ShippingDetails) <= 255, "shipping_details", "must not be more than 255 characters long")
}

func (p ProductModel) Update(prod *Product) error {
	query := `
    UPDATE products 
    SET product_name = $1, description = $2, price = $3, category = $4, materials_used = $5, shipping_details=$6, seller_id = $7
    WHERE product_id=$8
    `
	args := []interface{}{prod.Name, prod.Description, prod.Price, prod.Category, prod.MaterialsUsed, prod.ShippingDetails, prod.SellerID, prod.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, args...)
	return err

}

func (p ProductModel) Delete(id int) error {
	query := `
        DELETE FROM products 
        WHERE product_id=$1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, id)
	return err

}

func (p ProductModel) Insert(prod *Product) error {
	query := `
		INSERT INTO products (product_name, description, price, category, materials_used, shipping_details, seller_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING product_id
		`
	args := []interface{}{prod.Name, prod.Description, prod.Price, prod.Category, prod.MaterialsUsed, prod.ShippingDetails, prod.SellerID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&prod.ID)

}

func (p ProductModel) Get(id int) (*Product, error) {
	query := `
    SELECT product_id, product_name, description, price, category, materials_used, shipping_details, seller_id
    FROM products
    WHERE product_id=$1
    `
	var prod Product
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Category, &prod.MaterialsUsed, &prod.ShippingDetails, &prod.SellerID)
	if err != nil {
		return nil, err
	}
	return &prod, nil
}

func (p ProductModel) GetAll(category string, filters Filters) ([]*Product, Metadata, error) {
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), product_id, seller_id, product_name, description, price, category, materials_used, shipping_details
        FROM products
        WHERE category = $1
        ORDER BY %s %s, product_id
        LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{category, filters.limit(), filters.offset()}

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	products := []*Product{}

	for rows.Next() {
		var product Product
		err := rows.Scan(
			&totalRecords,
			&product.ID,
			&product.SellerID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Category,
			&product.MaterialsUsed,
			&product.ShippingDetails,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		products = append(products, &product)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return products, metadata, nil
}

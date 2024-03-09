package model 
/// we define API’s model
import (
    
	"database/sql"
	"log"
    "context"
    "time"
)

type Product struct {
    ID             int     `json:"id"`
    SellerID       int     `json:"seller_id"`
    Name           string  `json:"name"`
    Description    string  `json:"description"`
    Price          float64 `json:"price"`
    Category       string  `json:"category"`
    MaterialsUsed  string  `json:"materials_used"`
    ShippingDetails string `json:"shipping_details"`
}

type ProductModel struct { 
    DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (p ProductModel) Update (prod *Product) error { 
    query:=`
    UPDATE products 
    SET name = $1, description = $2, price = $3, category=$4, materials_used = $5, shipping_details=$6
    WHERE id=$7
    `
	args := []interface{}{prod.Name, prod.Description, prod.Price, prod.Category, prod.MaterialsUsed, prod.ShippingDetails}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return  p.DB.QueryRowContext(ctx, query, args...).Scan(&prod.ID)

}

func (p ProductModel) Delete(id int) error { 
    query:=`
        DELETE FROM products 
        WHERE id=$1
    `
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

    _, err := p.DB.ExecContext(ctx, query, id)
    return err

}

func (p ProductModel) Insert (prod *Product) error {
	query := `
		INSERT INTO products (seller_id, name, description, price, category, materials_used, shipping_details) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
		`
	args := []interface{}{prod.SellerID, prod.Name, prod.Description, prod.Price, prod.Category, prod.MaterialsUsed, prod.ShippingDetails}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return  p.DB.QueryRowContext(ctx, query, args...).Scan(&prod.ID, prod.SellerID)

}


func (p ProductModel) Get (id int) (*Product, error) { 
    query:= `
    SELECT id, seller_id, name, description, price, category, materials_used, shipping_details
    FROM products
    WHERE id=$1
    `
    var prod Product
    ctx, cancel:= context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&prod.ID, &prod.SellerID, &prod.Name, &prod.Description, &prod.Price, &prod.Category, &prod.MaterialsUsed, &prod.ShippingDetails)
	if err != nil {
		return nil, err
	}
	return &prod, nil
}



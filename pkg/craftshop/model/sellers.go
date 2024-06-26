package model

import (
	"context"
	"database/sql"
	"fmt"

	// "fmt"
	"log"
	"time"

	"github.com/ainelnazaraly/CraftShop/pkg/craftshop/validator"
)

type Seller struct {
	SellerID   int    `json:"seller_id"`
	SellerName string `json:"seller_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Location   string `json:"location"`
	DateJoined string `json:"date_joined"`
}
type SellerModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (s SellerModel) Insert(seller *Seller) error {
	query := `
	INSERT INTO sellers (seller_name, email, password, location)
	VALUES($1, $2, $3, $4) 
	RETURNING seller_id, date_joined
	`
	args := []interface{}{seller.SellerName, seller.Email, seller.Password, seller.Location}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&seller.SellerID, &seller.DateJoined)
}
func (s SellerModel) Get(id int) (*Seller, error) {
	query := `
        SELECT seller_id, seller_name, email, password, location
        FROM sellers
        WHERE seller_id = $1
    `
	var seller Seller
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&seller.SellerID, &seller.SellerName, &seller.Email, &seller.Password, &seller.Location)
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

func (s SellerModel) Update(seller *Seller) error {
	query := `
        UPDATE Sellers
        SET seller_name = $1, email = $2, password = $3, location = $4
        WHERE seller_id = $5
        RETURNING date_joined
    `
	args := []interface{}{seller.SellerName, seller.Email, seller.Password, seller.Location, seller.SellerID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&seller.DateJoined)
}

func (s SellerModel) Delete(id int) error {
	query := `
        DELETE FROM sellers
        WHERE seller_id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, query, id)
	return err
}

func (s SellerModel) GetAll(location string, filters Filters) ([]*Seller, Metadata, error) {
    query := fmt.Sprintf(`
        SELECT count(*) OVER(), seller_id, seller_name, email, password, location, date_joined
        FROM sellers
        WHERE location = $1
        ORDER BY %s %s, seller_id
        LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    args := []interface{}{location, filters.limit(), filters.offset()}

    rows, err := s.DB.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, Metadata{}, err
    }
    defer rows.Close()

    totalRecords := 0
    sellers := []*Seller{}

    for rows.Next() {
        var seller Seller
        err := rows.Scan(
            &totalRecords,
            &seller.SellerID,
            &seller.SellerName,
            &seller.Email,
            &seller.Password,
            &seller.Location,
            &seller.DateJoined,
        )
        if err != nil {
            return nil, Metadata{}, err
        }
        sellers = append(sellers, &seller)
    }
    if err = rows.Err(); err != nil {
        return nil, Metadata{}, err
    }

    metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
    return sellers, metadata, nil
}



func ValidateSeller(v *validator.Validator, seller *Seller) {
	// Check if the seller name field is empty.
	v.Check(seller.SellerName != "", "seller_name", "must be provided")
	// Check if the seller name field is not more than 100 characters.
	v.Check(len(seller.SellerName) <= 100, "seller_name", "must not be more than 100 characters long")
	// Check if the email field is empty.
	v.Check(seller.Email != "", "email", "must be provided")
	// Check if the email field is a valid email format.
	v.Check(seller.Email != "", "email", "must be a valid email address")
	// Check if the location field is empty.
	v.Check(seller.Location != "", "location", "must be provided")
	// Check if the location field is not more than 100 characters.
	v.Check(len(seller.Location) <= 100, "location", "must not be more than 100 characters long")
	v.Check(seller.Password != "", "password", "must be provided")
	// Check if the password field is at least 6 characters long.
	v.Check(len(seller.Password) >= 6, "password", "must be at least 6 characters long")
}

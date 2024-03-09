package model

import (
	"database/sql"
	"log"
    "context"
    "time"
)

type Seller struct { 
	SellerID	int 	`json:"seller_id"`
	SellerName	string 	`json:"seller_name"`
	Email 		string 	`json:"email"`
	Password 	string 	`json:"password"`
	Location 	string 	`json:"location"`
	DateJoined  string 	`json:"date_joined"`

} 
type SellerModel struct { 
    DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (s SellerModel) Insert (seller *Seller) error{ 
	query:=`
	INSERT INTO sellers (seller_name, email, password, location)
	VALUES($1, $2, $3) 
	RETURNING seller_id, date_joined
	`
	args:= []interface{}{seller.SellerName, seller.Email, seller.Password, seller.Location}
	ctx, cancel:=context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&seller.SellerID, &seller.DateJoined)
}

func (s SellerModel) Get (id int) (*Seller, error){ 
	query:=`
	SELECT seller_name, email, password, location
	FROM selllers 
	WHERE id=$1
	`
	var seller Seller 
	ctx, cancel:=context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row:= s.DB.QueryRowContext(ctx, query, id)
	err:=row.Scan(&seller.SellerName, &seller.Email, &seller.Password, &seller.Location)
	if err!=nil{ 
		return nil, err
	}
	return &seller, nil
}

func (s SellerModel) Update(seller *Seller) error{ 
	query:=`
	UPDATE sellers 
	SET seller_name =$1, email=$2, password=$3, location=$4
	WHERE id=$5
	RETURNING date_joined
	`
	args:= []interface{}{seller.SellerName, seller.Email, seller.Password, seller.Location}
	ctx, cancel :=context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&seller.DateJoined)
}

func (s SellerModel) Delete (id int) error{ 
	query:=`
	DELETE FROM sellers 
	WHERE id=$1
	`
	ctx, cancel:=context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err:=s.DB.ExecContext(ctx, query, id) 
	return err
}
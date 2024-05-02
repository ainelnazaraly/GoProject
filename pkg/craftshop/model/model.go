package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

type Models struct { 
	Products 	ProductModel
	Sellers 	SellerModel
}

var (
	// ErrRecordNotFound is returned when a record doesn't exist in database.
	ErrRecordNotFound = errors.New("record not found")

	// ErrEditConflict is returned when a there is a data race, and we have an edit conflict.
	ErrEditConflict = errors.New("edit conflict")
)

func NewModel(db *sql.DB) Models { 
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Products: ProductModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Sellers: SellerModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
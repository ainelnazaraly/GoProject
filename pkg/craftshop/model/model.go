package model

import (
	"database/sql"
	"log"
	"os"
)

type Models struct { 
	Products 	ProductModel
	Sellers 	SellerModel
}

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
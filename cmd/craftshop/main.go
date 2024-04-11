package main

import (
	"database/sql"
	"flag"
	"os"
	"net/http"


	"github.com/ainelnazaraly/CraftShop/pkg/craftshop/model"
	"github.com/ainelnazaraly/CraftShop/pkg/jsonlog"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
	logger *log.Logger
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:12345@localhost:5433/craftshop?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &application{
		config: cfg,
		models: model.NewModel(db),
		logger: logger
	}

	app.run()
}
func (app *application) run() {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Product Routes
	v1.HandleFunc("/products", app.createProductHandler).Methods("POST")
	v1.HandleFunc("/products/{productId:[0-9]+}", app.getProductHandler).Methods("GET")
	v1.HandleFunc("/products/{productId:[0-9]+}", app.updateProductHandler).Methods("PUT")
	v1.HandleFunc("/products/{productId:[0-9]+}", app.deleteProductHandler).Methods("DELETE")

// Create a new seller
v1.HandleFunc("/sellers", app.createSellerHandler).Methods("POST")
// Retrieve a seller by sellerName
v1.HandleFunc("/sellers/{sellerName}", app.getSellerHandler).Methods("GET")
// Update a seller's information by sellerName
v1.HandleFunc("/sellers/{sellerName}", app.updateSellerHandler).Methods("PUT")
// Delete a seller by sellerName
v1.HandleFunc("/sellers/{sellerName}", app.deleteSellerHandler).Methods("DELETE")


	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
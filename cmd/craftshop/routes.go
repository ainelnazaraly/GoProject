package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	log.Println("Starting API server")
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	//user

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Product Routes
	v1.HandleFunc("/products", app.createProductHandler).Methods("POST")
	v1.HandleFunc("/products/{product_id:[0-9]+}", app.getProductHandler).Methods("GET")
	v1.HandleFunc("/products/{product_id:[0-9]+}", app.updateProductHandler).Methods("PUT")
	v1.HandleFunc("/products/{product_id:[0-9]+}", app.requirePermissions("products:write", app.deleteProductHandler)).Methods("DELETE")

	// Create a new seller
	v1.HandleFunc("/sellers", app.requirePermissions("sellers:write", app.createSellerHandler)).Methods("POST")
	// Retrieve a seller by sellerName
	v1.HandleFunc("/sellers/{seller_id:[0-9]+}", app.requirePermissions("sellers:read", app.getSellerHandler)).Methods("GET")
	// Update a seller's information by sellerName
	v1.HandleFunc("/sellers/{seller_id:[0-9]+}", app.requirePermissions("sellers:write", app.updateSellerHandler)).Methods("PUT")
	// Delete a seller by sellerName
	v1.HandleFunc("/sellers/{seller_id:[0-9]+}", app.requirePermissions("sellers:write", app.deleteSellerHandler)).Methods("DELETE")

	// Retrieve a list of sellers
	v1.HandleFunc("/products", app.requireActivatedUser(app.listProductsHandler)).Methods("GET")

	v1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	v1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	v1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")
	return app.authenticate(r)
}

package main

import (
	"net/http"
	"strconv"

	"github.com/ainelnazaraly/CraftShop/pkg/craftshop/model"
	"github.com/gorilla/mux"
)

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name            string  `json:"name"`
		Description     string  `json:"description"`
		Price           float64 `json:"price"`
		Category        string  `json:"category"`
		MaterialsUsed   string  `json:"materials_used"`
		ShippingDetails string  `json:"shipping_details"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	product := &model.Product{
		Name:            input.Name,
		Description:     input.Description,
		Price:           input.Price,
		Category:        input.Category,
		MaterialsUsed:   input.MaterialsUsed,
		ShippingDetails: input.ShippingDetails,
	}

	err = app.models.Products.Insert(product)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, product)
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := app.models.Products.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 not found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, product)
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := app.models.Products.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 not found")
		return
	}

	var input struct {
		Name            *string  `json:"name"`
		Description     *string  `json:"description"`
		Price           *float64 `json:"price"`
		Category        *string  `json:"category"`
		MaterialsUsed   *string  `json:"materials_used"`
		ShippingDetails *string  `json:"shipping_details"`
	}
	err = app.readJSON(w, r, &input)

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		product.Name = *input.Name
	}

	if input.Description != nil {
		product.Description = *input.Description
	}

	if input.Price != nil {
		product.Price = *input.Price
	}

	if input.Category != nil {
		product.Category = *input.Category
	}

	if input.MaterialsUsed != nil {
		product.MaterialsUsed = *input.MaterialsUsed
	}

	if input.ShippingDetails != nil {
		product.ShippingDetails = *input.ShippingDetails
	}

	err = app.models.Products.Update(product)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, product)
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = app.models.Products.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

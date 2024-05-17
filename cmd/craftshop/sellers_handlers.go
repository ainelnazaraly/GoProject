package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ainelnazaraly/CraftShop/pkg/craftshop/model"
	"github.com/ainelnazaraly/CraftShop/pkg/craftshop/validator"
	// "github.com/ainelnazaraly/CraftShop/pkg/craftshop/validator"
	"github.com/gorilla/mux"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}
func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createSellerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SellerName string `json:"seller_name"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		Location   string `json:"location"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	seller := &model.Seller{
		SellerName: input.SellerName,
		Email:      input.Email,
		Password:   input.Password,
		Location:   input.Location,
	}

	err = app.models.Sellers.Insert(seller)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server error")
		return
	}
	app.respondWithJSON(w, http.StatusCreated, seller)
}

func (app *application) getSellerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["seller_id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid seller ID")
		return
	}

	seller, err := app.models.Sellers.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 not found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, seller)
}

func (app *application) updateSellerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["seller_id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid seller ID")
		return
	}

	seller, err := app.models.Sellers.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 not found")
		return
	}

	var input struct {
		SellerName *string `json:"seller_name"`
		Email      *string `json:"email"`
		Password   *string `json:"password"`
		Location   *string `json:"location"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.SellerName != nil {
		seller.SellerName = *input.SellerName
	}

	if input.Email != nil {
		seller.Email = *input.Email
	}

	if input.Password != nil {
		seller.Password = *input.Password
	}

	if input.Location != nil {
		seller.Location = *input.Location
	}

	err = app.models.Sellers.Update(seller)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, seller)
}

func (app *application) deleteSellerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["seller_id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid seller ID")
		return
	}

	err = app.models.Sellers.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) listSellersHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Location string
        model.Filters
    }
    v := validator.New()
    qs := r.URL.Query()

    input.Location = app.readString(qs, "location", "")
    input.Filters.Page = app.readInt(qs, "page", 1, v)
    input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
    input.Filters.Sort = app.readString(qs, "sort", "seller_id")

    input.Filters.SortSafeList = []string{"seller_id", "seller_name", "location", "-seller_id", "-seller_name", "-location"}

    if model.ValidateFilters(v, input.Filters); !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }

    sellers, metadata, err := app.models.Sellers.GetAll(input.Location, input.Filters)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, http.StatusOK, envelope{"sellers": sellers, "metadata": metadata}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}

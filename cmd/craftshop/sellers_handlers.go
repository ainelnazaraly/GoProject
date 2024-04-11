package main

import (
	"encoding/json"
	"net/http"

	"github.com/ainelnazaraly/CraftShop/pkg/craftshop/model"
	"github.com/ainelnazaraly/CraftShop/pkg/craftshop/validator"
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) createSellerHandler(w http.ResponseWriter, r *http.Request){ 
	var input struct { 
		SellerName	string 	`json:"seller_name"`
		Email 		string 	`json:"email"`
		Password 	string 	`json:"password"`
		Location 	string 	`json:"location"`
	}
	err:=app.readJSON(w, r, &input)
	if err!= nil{ 
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	seller:=&model.Seller{ 
		SellerName:  input.SellerName,
		Email: input.Email,
		Password: input.Password,
		Location: input.Location,
	}

	err =app.models.Sellers.Insert(seller)
	if err!=nil{ 
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server error")
		return
	}
	app.respondWithJSON(w, http.StatusCreated, seller)
}

func (app *application) getSellerHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["sellerName"]

    if name == "" {
        app.respondWithError(w, http.StatusBadRequest, "Invalid seller name")
        return
    }

    seller, err := app.models.Sellers.Get(name)
    if err != nil {
        app.respondWithError(w, http.StatusNotFound, "Seller not found")
        return
    }

    app.respondWithJSON(w, http.StatusOK, seller)
}

func (app *application) updateSellerHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["sellerName"]

    if name == "" {
        app.respondWithError(w, http.StatusBadRequest, "Invalid seller name")
        return
    }

    seller, err := app.models.Sellers.Get(name)
    if err != nil {
        app.respondWithError(w, http.StatusNotFound, "Seller not found")
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
    name := vars["sellerName"]

    if name == "" {
        app.respondWithError(w, http.StatusBadRequest, "Invalid seller name")
        return
    }

    err := app.models.Sellers.Delete(name)
    if err != nil {
        app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
        return
    }

    app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) getSellersList(w http.ResponseWriter, r *http.Request){ 
    var input struct { 	
	    SellerName	string 		
	    Location 	[]string
        Page        int 
        PageSize    int
        Sort        string
    }

    v := validator.New()

    qs := r.URL.Query()

    input.SellerName = app.readString(qs, "seller_name", "")
    input.Location= app.readCSV(qs, "location", []string{})

    input.Page =app.readInt(qs, "page", 1, v)
    input.PageSize = app.readInt(qs, "page_size", 20, v)

    input.Sort = app.readString(qs, "sort", "id")
    
    if !v.Valid(){ 
        app.fa
    }

}

package filler

import "github.com/ainelnazaraly/CraftShop/pkg/craftshop/model"

func PopulateDatabase(models model.Models) error {
	for _, product := range products {
		err := models.Products.Insert(&product)
		if err != nil {
			return err
		}
	}
	return nil
}

var products = []model.Product{
	{Name: "Organic Cotton T-Shirt", Description: "Sustainably sourced cotton t-shirt", Price: 25.99, Category: "Clothing", MaterialsUsed: "Organic Cotton", ShippingDetails: "Standard Shipping", SellerID: 4},
	{Name: "Handmade Ceramic Plant Pot", Description: "Artisan crafted ceramic plant pot", Price: 19.99, Category: "Home Decor", MaterialsUsed: "Ceramic", ShippingDetails: "Free Shipping", SellerID: 3},
	// Add more products here
}


# CraftShop API

The CraftShop API provides endpoints for managing products and the corresponding sellers in a craftshop application. 

# DataBase Schema 
![DataBaseScheme](https://github.com/ainelnazaraly/GoProject/blob/fd71f93a5ca0150c22256e58fa329c4dcf6843e1/dbscheme.png)

# DB structure 
table sellers {  
  seller_id integer [primary key]


  seller_name varchar
  
  email varchar
  
  password varchar
  
  location varchar 
  
  date_joined timestamp
  
}

table products { 
  product_id int [primary key]
  
  seller_id int 
  
  product_name varchar
  
  description varchar 
  
  price varchar 
  
  category varchar 
  
  materials_used varchar 
  
  shipping_details varchar 
}


Ref: "sellers"."seller_id" < "products"."seller_id" 

# Available Endpoints

## Product Routes:

POST /api/v1/products: Creates a new product.

GET /api/v1/products/{productId}: Retrieves details of a product by its ID.

PUT /api/v1/products/{productId}: Updates details of a product by its ID.

DELETE /api/v1/products/{productId}: Deletes a product by its ID.

## Seller Routes:

POST /api/v1/sellers: Creates a new seller.

GET /api/v1/sellers/{sellerId}: Retrieves details of a seller by their ID.

PUT /api/v1/sellers/{sellerId}: Updates information of a seller by their ID.

DELETE /api/v1/sellers/{sellerId}: Deletes a seller by their ID.

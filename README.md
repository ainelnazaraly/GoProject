#CraftShop API
The CraftShop API provides endpoints for managing products and the corresponding sellers in a craftshop application 

##DataBase Schema 
![DataBaseScheme](https://github.com/ainelnazaraly/GoProject.git/dbscheme.png)

##DB structure 
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

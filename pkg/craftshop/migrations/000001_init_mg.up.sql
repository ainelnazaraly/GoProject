
CREATE TABLE  IF NOT EXISTS Sellers (
    seller_id bigserial PRIMARY KEY,
    seller_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    location VARCHAR(255),
    date_joined TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Products (
    product_id bigserial PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(100),
    materials_used TEXT,
    shipping_details TEXT,
    seller_id bigserial references Sellers(seller_id)
);

select * from products;

select * from sellers;

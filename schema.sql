CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    image_url TEXT
);

CREATE TABLE suppliers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(10) DEFAULT '₽',
    url TEXT,
    notes TEXT,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

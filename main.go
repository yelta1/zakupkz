package main

import (
    "database/sql"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gin-contrib/cors"
)

type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Category    string  `json:"category"`
    ImageURL    string  `json:"image_url"`
}

type Supplier struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    Currency string  `json:"currency"`
    URL      string  `json:"url"`
    Notes    string  `json:"notes"`
}

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("mysql", "username:password@tcp(localhost:3306)/analka_db?parseTime=true")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    r := gin.Default()
    r.Use(cors.Default()) // <--- вот эта строка разрешает CORS

    r.GET("/products", getProducts)
    r.POST("/products", addProduct)
    r.GET("/products/:id/suppliers", getSuppliers)
    r.POST("/products/:id/suppliers", addSupplier)
    r.DELETE("/products/:id", deleteProduct)

    r.Run(":8080")
}

func getProducts(c *gin.Context) {
    rows, err := db.Query("SELECT id, name, description, category, image_url FROM products")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var p Product
        rows.Scan(&p.ID, &p.Name, &p.Description, &p.Category, &p.ImageURL)
        products = append(products, p)
    }
    c.JSON(http.StatusOK, products)
}

func addProduct(c *gin.Context) {
    var p Product
    if err := c.ShouldBindJSON(&p); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    res, err := db.Exec("INSERT INTO products (name, description, category, image_url) VALUES (?, ?, ?, ?)",
        p.Name, p.Description, p.Category, p.ImageURL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    id, _ := res.LastInsertId()
    p.ID = int(id)
    c.JSON(http.StatusOK, p)
}

func getSuppliers(c *gin.Context) {
    productID := c.Param("id")
    rows, err := db.Query("SELECT id, name, price, currency, url, notes FROM suppliers WHERE product_id = ?", productID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var suppliers []Supplier
    for rows.Next() {
        var s Supplier
        rows.Scan(&s.ID, &s.Name, &s.Price, &s.Currency, &s.URL, &s.Notes)
        suppliers = append(suppliers, s)
    }
    c.JSON(http.StatusOK, suppliers)
}

func addSupplier(c *gin.Context) {
    productID := c.Param("id")
    var s Supplier
    if err := c.ShouldBindJSON(&s); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    _, err := db.Exec("INSERT INTO suppliers (product_id, name, price, currency, url, notes) VALUES (?, ?, ?, ?, ?, ?)",
        productID, s.Name, s.Price, s.Currency, s.URL, s.Notes)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func deleteProduct(c *gin.Context) {
    id := c.Param("id")
    _, err := db.Exec("DELETE FROM products WHERE id = ?", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "deleted"})
} 

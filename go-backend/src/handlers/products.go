package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"go-backend/src/middleware"
)

type Product struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	PriceInDollar string `json:"price_in_dollar"`
	ThumbnailUrl  string `json:"thumbnail_url"`
}

func GetProducts(c *gin.Context) {
	dbService, exists := c.Get("DbService")
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"error": "Database service not found"})
		return
	}
	db := dbService.(*middleware.DatabaseService)
	var products []Product

	rows, err := db.Pool.Query(context.Background(), "SELECT * FROM Products")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database query error: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server Error"})
		return
	}

	for rows.Next() {
		var product Product

		err := rows.Scan(&product.Name, &product.Description, &product.PriceInDollar, &product.ThumbnailUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Database query error: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server Error"})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

func PostProduct(c *gin.Context) {
	dbService, exists := c.Get("DbService")
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"error": "Database service not found"})
		return
	}
	db := dbService.(*middleware.DatabaseService)
	var product Product

	buf := []byte{}
	n, err := c.Request.Body.Read(buf)
	if err != nil || n < 0 {
		fmt.Fprintf(os.Stderr, "Bad request body: %v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
	if json.Unmarshal(buf, &product) != nil {
		fmt.Fprintf(os.Stderr, "Couldn't unmarshal, bad request body: %v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	println(product.Name)
	println(product.ThumbnailUrl)

	conn, err := db.Pool.Exec(
		context.Background(),
		"INSERT INTO Products (name, description, price_in_dollar, thumbnail_url) VALUES ($1, $2, $3, $4)",
		product.Name,
		product.Description,
		product.PriceInDollar,
		product.ThumbnailUrl,
	)

	if err != nil || conn.RowsAffected() <= 0 {
		fmt.Fprintf(os.Stderr, "Database query error: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server Error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created"})
}

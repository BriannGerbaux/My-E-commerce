package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

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

func ListProducts(c *gin.Context) {
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

func GetProduct(c *gin.Context) {
	dbService, exists := c.Get("DbService")
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"error": "Database service not found"})
		return
	}
	db := dbService.(*middleware.DatabaseService)

	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	var product Product

	row := db.Pool.QueryRow(context.Background(), "SELECT * FROM Products WHERE id = $1", productId)
	err = row.Scan(&product.Name, &product.Description, &product.PriceInDollar, &product.ThumbnailUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database query error: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server Error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func PostProduct(c *gin.Context) {
	dbService, exists := c.Get("DbService")
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"error": "Database service not found"})
		return
	}
	db := dbService.(*middleware.DatabaseService)
	var product Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

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

func UpdateProduct(c *gin.Context) {
	dbService, exists := c.Get("DbService")
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"error": "Database service not found"})
		return
	}
	db := dbService.(*middleware.DatabaseService)
	var bodyProduct Product

	if err := c.ShouldBindJSON(&bodyProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	conn, err := db.Pool.Exec(
		context.Background(),
		"UPDATE Products SET name = $1, description = $2, price_in_dollar = $3, thumbnail_url = $4 WHERE id = $5",
		bodyProduct.Name,
		bodyProduct.Description,
		bodyProduct.PriceInDollar,
		bodyProduct.ThumbnailUrl,
		productId,
	)

	if err != nil || conn.RowsAffected() <= 0 {
		fmt.Fprintf(os.Stderr, "Database query error: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server Error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created"})
}

func DeleteProduct(c *gin.Context) {
	dbService, exists := c.Get("DbService")
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"error": "Database service not found"})
		return
	}
	db := dbService.(*middleware.DatabaseService)

	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	conn, err := db.Pool.Exec(context.Background(), "DELETE FROM Products WHERE id = $1", productId)

	if err != nil || conn.RowsAffected() <= 0 {
		fmt.Fprintf(os.Stderr, "Database query error: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

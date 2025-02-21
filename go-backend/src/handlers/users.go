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

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetUserById(c *gin.Context) {
	dbService, exists := c.Get("DbService")
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"error": "Database service not found"})
		return
	}
	db := dbService.(*middleware.DatabaseService)
	var id int

	id, err := strconv.Atoi(c.GetString("user_id"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Id %s cannot be converted to integer: %v\n", c.GetString("user_id"), err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	var name string
	var email string

	err = db.Pool.QueryRow(context.Background(), "SELECT name, email FROM Users WHERE id=$1", id).Scan(&name, &email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database query error: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		name:  name,
		email: email,
	})
}

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

func ListUsers(c *gin.Context) {
	dbService, exists := c.Get("DbService")
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"error": "Database service not found"})
		return
	}
	db := dbService.(*middleware.DatabaseService)

	var users []User

	rows, err := db.Pool.Query(context.Background(), "SELECT id, name, email FROM Users")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database query error: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server Error"})
		return
	}

	for rows.Next() {
		var user User

		err = rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Rows scan error: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server Error"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func GetUserById(c *gin.Context) {
	dbService, exists := c.Get("DbService")
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"error": "Database service not found"})
		return
	}
	db := dbService.(*middleware.DatabaseService)
	var id int

	userIdAsAny, userIdExists := c.Get("user_id")
	_, adminIdExists := c.Get("admin_id")
	// If is from User, the user ID requested to get must be the same as the ID of the user that made the request
	if userIdExists {
		requestedUserId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Wrong request user ID: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}
		loggedInUserId, err := strconv.Atoi(userIdAsAny.(string))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Wrong logged user ID: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}
		if loggedInUserId != requestedUserId {
			fmt.Fprintln(os.Stderr, "Logged in user ID and requested user ID does not match. A user can only fetch it's own account information")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}
		id = requestedUserId
	} else if adminIdExists { // If is from admin, just use the id in parameters. No restriction
		convertedId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Wrong request user ID: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}
		id = convertedId
	} else {
		fmt.Fprintln(os.Stderr, "No admin or user ID found")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	var userId string
	var name string
	var email string

	err := db.Pool.QueryRow(context.Background(), "SELECT id, name, email FROM Users WHERE id=$1", id).Scan(&userId, &name, &email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database query error: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    userId,
		"name":  name,
		"email": email,
	})
}

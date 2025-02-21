package middleware

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseService struct {
	Pool *pgxpool.Pool
}

func DbConnection() DatabaseService {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return DatabaseService{Pool: dbpool}
}

func InjectDatabaseService(dbService *DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DbService", dbService)
		c.Next()
	}
}

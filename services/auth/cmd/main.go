package main

import (
	"context"
	"go-template/gin_sqlc_setup/services/auth/db"
	"go-template/gin_sqlc_setup/services/auth/handlers"
	"go-template/gin_sqlc_setup/services/auth/middleware"
	"go-template/gin_sqlc_setup/services/auth/repository"
	"go-template/gin_sqlc_setup/services/auth/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func openDBConnection(dsn string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func closeDBConnection(pool *pgxpool.Pool) {
	pool.Close()
}
func main() {
	dsn := "postgres://sprice:sprice@localhost:5432/template-auth?sslmode=disable"
	dbConn, err := openDBConnection(dsn)
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(dbConn)
	queries := db.New(dbConn)

	userRepo := repository.NewUserRepository(queries)
	authService := service.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()
	r.POST("/login", authHandler.Login)
	r.POST("/register", authHandler.Register)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/protected", authHandler.Profile)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

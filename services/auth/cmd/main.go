package main

import (
	"context"
	"fmt"
	"go-template/gin_sqlc_setup/services/auth/config"
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

	conf := config.LoadConfig()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
	)
	dbConn, err := openDBConnection(dsn)
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(dbConn)
	queries := db.New(dbConn)

	userRepo := repository.NewUserRepository(queries)
	authService := service.NewAuthService(userRepo, conf)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()
	r.POST("/login", authHandler.Login)
	r.POST("/register", authHandler.Register)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(conf))
	protected.GET("/profile", authHandler.Profile)

	if err := r.Run(fmt.Sprintf("%s:%s", conf.Host, conf.Port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

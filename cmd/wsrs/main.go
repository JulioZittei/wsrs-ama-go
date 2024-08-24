package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/JulioZittei/wsrs-ama-go/docs"
	"github.com/JulioZittei/wsrs-ama-go/internal/app"
	"github.com/JulioZittei/wsrs-ama-go/internal/store/pgstore"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// @title Ask Me Anything API
// @version 1.0
// @description This is a API documentation for Ask Me Anything
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("WSRS_DATABASE_USER"), os.Getenv("WSRS_DATABASE_PASSWORD"), os.Getenv("WSRS_DATABASE_HOST"), os.Getenv("WSRS_DATABASE_PORT"), os.Getenv("WSRS_DATABASE_NAME")))

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	app := app.NewApplication(pgstore.New(pool))
	app.Init()

	server := &http.Server{
		Addr:    ":8080",
		Handler: app.GetHandler(),
	}

	go func() {
		fmt.Println("Starting Server")
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)

	fmt.Println("Started")
	<-quit

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("ServerStopping...")

	if err := server.Shutdown(context); err != nil {
		panic(err)
	}

	fmt.Println("Server Stopped")
}

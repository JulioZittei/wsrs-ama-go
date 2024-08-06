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

	"github.com/JulioZittei/wsrs-ama-go/internal/api"
	"github.com/JulioZittei/wsrs-ama-go/internal/store/pgstore"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

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

	handler := api.NewHandler(pgstore.New(pool))

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
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

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vykiy/house-service/internal/app"
	"github.com/Vykiy/house-service/internal/config"
	"github.com/Vykiy/house-service/internal/repository"
	"github.com/Vykiy/house-service/internal/router"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	config := config.NewConfig()

	db, err := sqlx.Connect("postgres", config.DBConnection)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	repo := repository.NewRepository(db)

	app := app.NewApp(repo)

	jwtIssuer := router.NewJWTIssuer(config.JWTSecret)

	router := router.NewRouter(app, jwtIssuer)

	server := &http.Server{
		Addr:    config.ServerAddress,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()
	log.Printf("Server is ready to handle requests at %s", server.Addr)

	<-quit
	log.Printf("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	log.Printf("Server stopped")
}

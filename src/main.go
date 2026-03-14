package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"supportflow/core"
	"supportflow/db/postgre"
	dbRedis "supportflow/db/redis"
	"supportflow/routes"
	"supportflow/services/ai"
)

func main() {
	configPath := "../config/local-application.properties"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		configPath = p
	}
	if err := core.LoadConfig(configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := postgre.Init(ctx); err != nil {
		log.Fatalf("Postgres init failed: %v", err)
	}
	defer postgre.Close()

	if err := postgre.RunMigrations(ctx); err != nil {
		log.Printf("Migrations warning: %v", err)
	}

	ai.Init()

	if err := dbRedis.Init(ctx); err != nil {
		log.Printf("Redis init failed (continuing without): %v", err)
	}
	defer dbRedis.Close()

	router := mux.NewRouter()
	routes.Register(router, ctx)

	origins := strings.Split(core.GetString("cors.allowed_origins", "http://localhost:5173"), ",")
	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	port := core.GetString("service.port", "8080")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      c.Handler(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("SupportFlow AI starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}

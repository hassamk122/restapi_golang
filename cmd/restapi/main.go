package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hassamk122/restapi_golang/config"
	"github.com/hassamk122/restapi_golang/config/dbconfig"
	"github.com/hassamk122/restapi_golang/internal/handlers"
	"github.com/hassamk122/restapi_golang/internal/middlewares"
	"github.com/hassamk122/restapi_golang/internal/repo"
	"github.com/hassamk122/restapi_golang/internal/routes"
	"github.com/hassamk122/restapi_golang/internal/service"
	"github.com/hassamk122/restapi_golang/internal/store"
	"github.com/redis/go-redis/v9"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config : %v", err)
		os.Exit(1)
	}

	db := dbconfig.ConnectDB(config.DatabaseUrl)
	defer db.Close()

	rdb := dbconfig.ConnectRedis()
	defer func(rdb *redis.Client) {
		_ = rdb.Close()
	}(rdb)

	queries := store.New(db)

	userRepo := repo.NewUserRepo(queries)

	userService := service.NewUserService(db, userRepo, rdb)

	handler := handlers.NewHandler(userService)

	mux := http.NewServeMux()
	routes.SetupRoutes(mux, handler)

	loggedMux := middlewares.LoggingMiddleware(mux)

	serverAddr := fmt.Sprintf(":%s", config.ServerPort)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: loggedMux,
	}

	log.Printf("Server Started at Port %s\n", serverAddr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed %v", err)
	}
}

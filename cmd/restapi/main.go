package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hassamk122/restapi_golang/config"
	"github.com/hassamk122/restapi_golang/config/dbconfig"
	"github.com/hassamk122/restapi_golang/internal/handlers"
	"github.com/hassamk122/restapi_golang/internal/repo"
	"github.com/hassamk122/restapi_golang/internal/routes"
	"github.com/hassamk122/restapi_golang/internal/service"
	"github.com/hassamk122/restapi_golang/internal/store"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config : %v", err)
		os.Exit(1)
	}

	db := dbconfig.ConnectDB(config.DatabaseUrl)
	defer db.Close()

	queries := store.New(db)

	userRepo := repo.NewUserRepo(queries)
	userService := service.NewUserService(db, userRepo)
	handler := handlers.NewHandler(userRepo, userService)

	mux := http.NewServeMux()
	routes.SetupRoutes(mux, handler)

	serverAddr := fmt.Sprintf(":%s", config.ServerPort)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}

	log.Printf("Server Started at Port %s\n", serverAddr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed %v", err)
	}
}

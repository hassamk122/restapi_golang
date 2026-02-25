package routes

import (
	"net/http"

	"github.com/hassamk122/restapi_golang/internal/handlers"
	"github.com/hassamk122/restapi_golang/internal/utils"
)

func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	SetupHealthRoute(mux, handler)
	SetupUserRoutes(mux, handler)

	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		utils.ResondWithNotFound(res)
	})
}

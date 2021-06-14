package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	// this line is used by starport scaffolding # 1
)

const (
	MethodGet = "GET"
)

// RegisterRoutes registers verim-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 2
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)

}

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 3
	r.HandleFunc("/verim/nyms/{id}", getNymHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/verim/nyms", listNymHandler(clientCtx)).Methods("GET")

}

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 4
	r.HandleFunc("/verim/nyms", createNymHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/verim/nyms/{id}", updateNymHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/verim/nyms/{id}", deleteNymHandler(clientCtx)).Methods("POST")

}

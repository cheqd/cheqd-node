package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	// this line is used by starport scaffolding # 1
)

const (
	MethodGet = "GET"
)

// RegisterRoutes registers cheqd-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 2
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)

}

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 3
	r.HandleFunc("/cheqd/nyms/{id}", getNymHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/cheqd/nyms", listNymHandler(clientCtx)).Methods("GET")

}

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 4
	r.HandleFunc("/cheqd/nyms", createNymHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/cheqd/nyms/{id}", updateNymHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/cheqd/nyms/{id}", deleteNymHandler(clientCtx)).Methods("POST")

}

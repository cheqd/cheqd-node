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
	r.HandleFunc("/cheqd/credDefs/{id}", getCredDefHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/cheqd/credDefs", listCredDefHandler(clientCtx)).Methods("GET")

	r.HandleFunc("dock", getSchemaHandler(clientCtx)).Methods("GET")

	r.HandleFunc("/cheqd/dids/{id}", getDidHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/cheqd/dids", listDidHandler(clientCtx)).Methods("GET")
}

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 4
	r.HandleFunc("/cheqd/credDefs", createCredDefHandler(clientCtx)).Methods("POST")

	r.HandleFunc("/cheqd/schemata", createSchemaHandler(clientCtx)).Methods("POST")

	r.HandleFunc("/cheqd/dids", createDidHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/cheqd/dids/{id}", updateDidHandler(clientCtx)).Methods("POST")
}

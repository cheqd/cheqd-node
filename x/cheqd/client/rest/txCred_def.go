package rest

import (
	"net/http"
	"strconv"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

type createCred_defRequest struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	Creator        string       `json:"creator"`
	Schema_id      string       `json:"schema_id"`
	Tag            string       `json:"tag"`
	Signature_type string       `json:"signature_type"`
	Value          string       `json:"value"`
}

func createCred_defHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createCred_defRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		parsedSchema_id := req.Schema_id

		parsedTag := req.Tag

		parsedSignature_type := req.Signature_type

		parsedValue := req.Value

		msg := types.NewMsgCreateCred_def(
			req.Creator,
			parsedSchema_id,
			parsedTag,
			parsedSignature_type,
			parsedValue,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type updateCred_defRequest struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	Creator        string       `json:"creator"`
	Schema_id      string       `json:"schema_id"`
	Tag            string       `json:"tag"`
	Signature_type string       `json:"signature_type"`
	Value          string       `json:"value"`
}

func updateCred_defHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			return
		}

		var req updateCred_defRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err = sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		parsedSchema_id := req.Schema_id

		parsedTag := req.Tag

		parsedSignature_type := req.Signature_type

		parsedValue := req.Value

		msg := types.NewMsgUpdateCred_def(
			req.Creator,
			id,
			parsedSchema_id,
			parsedTag,
			parsedSignature_type,
			parsedValue,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type deleteCred_defRequest struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Creator string       `json:"creator"`
}

func deleteCred_defHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			return
		}

		var req deleteCred_defRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		_, err = sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgDeleteCred_def(
			req.Creator,
			id,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

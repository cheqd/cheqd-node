package rest

import (
	"net/http"
	"strconv"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

type DidService struct {
	Id              string `json:"id""`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

type VerificationMethod struct {
	Id                 string            `json:"id"`
	Type               string            `json:"type"`
	Controller         string            `json:"controller"`
	PublicKeyJwk       map[string]string `json:"publicKeyJwk"`
	PublicKeyMultibase string            `json:"publicKeyMultibase"`
}

type createDidRequest struct {
	BaseReq              rest.BaseReq          `json:"base_req"`
	Id                   string                `json:"id"`
	Controller           []string              `json:"controller"`
	VerificationMethod   []*VerificationMethod `json:"alias"`
	Authentication       []string              `json:"authentication"`
	AssertionMethod      []string              `json:"assertionMethod"`
	CapabilityInvocation []string              `json:"capabilityInvocation"`
	CapabilityDelegation []string              `json:"capabilityDelegation"`
	KeyAgreement         []string              `json:"keyAgreement"`
	AlsoKnownAs          []string              `json:"alsoKnownAs"`
	Service              []*DidService         `json:"service"`
}

func createDidHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createDidRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// TODO add verifificationMethod and Service
		msg := types.NewMsgCreateDid(
			req.Id,
			req.Controller,
			nil,
			req.Authentication,
			req.AssertionMethod,
			req.CapabilityInvocation,
			req.CapabilityDelegation,
			req.KeyAgreement,
			req.AlsoKnownAs,
			nil,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type updateDidRequest struct {
	BaseReq              rest.BaseReq          `json:"base_req"`
	Id                   string                `json:"id"`
	Controller           []string              `json:"controller"`
	VerificationMethod   []*VerificationMethod `json:"alias"`
	Authentication       []string              `json:"authentication"`
	AssertionMethod      []string              `json:"assertionMethod"`
	CapabilityInvocation []string              `json:"capabilityInvocation"`
	CapabilityDelegation []string              `json:"capabilityDelegation"`
	KeyAgreement         []string              `json:"keyAgreement"`
	AlsoKnownAs          []string              `json:"alsoKnownAs"`
	Service              []*DidService         `json:"service"`
}

func updateDidHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			return
		}

		var req updateDidRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// TODO add verifificationMethod and Service
		msg := types.NewMsgUpdateDid(
			req.Id,
			req.Controller,
			nil,
			req.Authentication,
			req.AssertionMethod,
			req.CapabilityInvocation,
			req.CapabilityDelegation,
			req.KeyAgreement,
			req.AlsoKnownAs,
			nil,
		)

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

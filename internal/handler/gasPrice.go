package handler

import (
	"encoding/json"
	"net/http"

	"interview-test/internal/client/infura"
)

type GetGasPriceHandler struct {
	client infura.Client
}

func NewGetGasPriceHandler(client infura.Client) http.Handler {
	return &GetGasPriceHandler{client: client}
}

type getGasPriceResponse struct {
	GasPrice string `json:"gasPrice"`
}

func (h GetGasPriceHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if wei, err := h.client.GetGasPrice(req.Context()); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	} else {
		_ = json.NewEncoder(rw).Encode(&getGasPriceResponse{GasPrice: wei})
	}
}

package api

import (
	"encoding/json"
	"net/http"
)

type AccountAPI struct {
	Service AccountService
}

type AccountService interface {
	CreateAccount(AccountInput) (any, error)
	ListAccounts() (any, error)
}

type AccountInput struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Currency string `json:"currency"`
}

func (a AccountAPI) Create(w http.ResponseWriter, r *http.Request) {
	var input AccountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if input.Name == "" || input.Type == "" || input.Currency == "" {
		http.Error(w, "name, type and currency are required", http.StatusBadRequest)
		return
	}
	result, err := a.Service.CreateAccount(input)
	if err != nil {
		http.Error(w, "account creation failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
}

func (a AccountAPI) List(w http.ResponseWriter, r *http.Request) {
	result, err := a.Service.ListAccounts()
	if err != nil {
		http.Error(w, "account lookup failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

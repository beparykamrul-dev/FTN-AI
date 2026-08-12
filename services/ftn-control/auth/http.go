package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

type SignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// SignInHandler is transport-only: credentials are verified by the auth service.
// The password is never persisted or logged here.
func SignInHandler(service *Service, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req SignInRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password are required", http.StatusBadRequest)
		return
	}

	// Salt lookup/storage is intentionally kept in the provisioning repository;
	// this handler does not invent or persist credentials.
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		http.Error(w, "authentication unavailable", http.StatusInternalServerError)
		return
	}
	result, err := service.Authenticate(r.Context(), req.Login, []byte(req.Password), salt)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:     "ftn_session",
		Value:    base64.RawURLEncoding.EncodeToString([]byte(result.Session.ID)),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"identity_id": result.IdentityID})
}

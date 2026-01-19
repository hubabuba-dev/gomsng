package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"msng/internal/service"

	"github.com/google/uuid"
)

type UserHandler struct {
	userService *service.UserService
	authService *service.AuthService
}

// Requests
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Responses
type RegisterResponse struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
}

type LoginRefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

type RefreshCookie struct {
	RefreshToken *http.Cookie
}

func CreateUserHandler(userService *service.UserService, authService *service.AuthService) *UserHandler {
	return &UserHandler{userService: userService,
		authService: authService}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request RegisterRequest
	var response RegisterResponse

	json_decoder := json.NewDecoder(r.Body)
	json_decoder.DisallowUnknownFields()
	if err := json_decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if request.Username == "" || request.Password == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}
	var err error

	response.UUID, response.Username, err = h.userService.RegisterUser(ctx, request.Username, request.Password)
	if err != nil {
		http.Error(w, "Error while creating user", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	var request LoginRequest
	var response LoginRefreshResponse
	var refresh_token string

	json_decoder := json.NewDecoder(r.Body)
	json_decoder.DisallowUnknownFields()
	if err := json_decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if request.Username == "" || request.Password == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	response.AccessToken, refresh_token, err = h.authService.LoginUser(ctx, request.Username, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBlocked):
			http.Error(w, "User is blocked", http.StatusForbidden)
		case errors.Is(err, service.ErrWrongCreds):
			http.Error(w, "Wrong credentials", http.StatusUnauthorized)
		case errors.Is(err, service.ErrTokenNotValid):
			http.Error(w, "Token not valid", http.StatusUnauthorized)
		default:
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
		}
		return
	}

	refresh_cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refresh_token,
		MaxAge:   3600,
		HttpOnly: true,
		Path:     "/auth",
		//Secure: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &refresh_cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) RefreshUserToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Refresh token start")
	ctx := r.Context()
	var err error

	var request RefreshCookie
	var response LoginRefreshResponse
	var refresh_token string

	request.RefreshToken, err = r.Cookie("refreshToken")
	if err != nil {
		log.Println("No refresh token fallback")
		http.Error(w, "No refresh token", http.StatusUnauthorized)
		return
	}

	refreshToken := request.RefreshToken.Value

	response.AccessToken, refresh_token, err = h.authService.RefreshUserToken(ctx, refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBlocked):
			http.Error(w, "User is blocked", http.StatusForbidden)
		case errors.Is(err, service.ErrTokenNotValid):
			http.Error(w, "Token not valid", http.StatusUnauthorized)
		default:
			log.Println(err)
			http.Error(w, "Error while creating token", http.StatusInternalServerError)
		}
		return
	}

	refresh_cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refresh_token,
		MaxAge:   3600,
		HttpOnly: true,
		Path:     "/auth",
		//Secure: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &refresh_cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	refresh_token, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, "Empty refresh token", http.StatusForbidden)
		return
	}

	refreshTokenValue := refresh_token.Value

	err = h.authService.LogoutUser(ctx, refreshTokenValue)
	if err != nil {
		http.Error(w, "Smth baaad", http.StatusBadGateway)
		return
	}
	refresh_cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		MaxAge:   -1,
		Path:     "/auth",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, refresh_cookie)
	w.WriteHeader(http.StatusNoContent)
}

package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"userService/internal/auth/tokenHandler"
	"userService/user/model"
	"userService/user/service"

	"github.com/go-chi/chi"
)

const authorization = "Authorization"

type UserService interface {
	GetUser(email string) (model.User, error)
	GetAll() ([]model.User, error)
	Update(u model.User) error
	Create(u model.User) error
	Delete(username string) error
	Register(user model.User) error
	AuthorizeEmail(user service.LoginRequest) (string, string, error)
	ParseToken(inputToken string, kind int) (service.UserClaims, error)
	GenerateTokens(email string) (string, string, error)
}

func NewHanadler(auth UserService) Handler {
	return Handler{auth: auth}
}

func (h Handler) NewApiRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Route("/refresh", func(r chi.Router) {
			r.Use(h.CheckRefresh)
			r.Post("/", h.Refresh)
		})
	})

	return r
}

func (a Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	val, ok := r.Context().Value(UserRequest{}).(service.UserClaims)
	if !ok {
		http.Error(w, "cant parse data", http.StatusBadRequest)
		return
	}

	user, err := a.auth.GetUser(val.Email)
	if err != nil {
		http.Error(w, "cant auth", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := a.auth.GenerateTokens(user.Email)
	if err != nil {
		http.Error(w, "cant auth", http.StatusBadRequest)
		return
	}

	rt := loginData{
		Message:      "success refresh",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	b, err := json.Marshal(rt)
	if err != nil {
		http.Error(w, "something very wrong", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(b)
	w.WriteHeader(http.StatusOK)
}

func (h Handler) CheckRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenRaw := r.Header.Get(authorization)
		tokenParts := strings.Split(tokenRaw, " ")
		if len(tokenParts) < 2 && tokenParts[0] != "Bearer" {
			http.Error(w, "wrong input data", http.StatusBadRequest)
			return
		}
		u, err := h.auth.ParseToken(tokenParts[1], tokenHandler.AccessToken)
		if err != nil && err.Error() == "unknown type of token" {
			http.Error(w, "unknown type of token", http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), UserRequest{}, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user service.LoginRequest
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cant read data", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(b, &user)
	if err != nil {
		http.Error(w, "cant parse data", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.auth.AuthorizeEmail(user)
	if err != nil {
		http.Error(w, "auth failed", http.StatusBadRequest)
		return
	}
	rt := loginData{
		Message:      "success login",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	b, err = json.Marshal(rt)
	if err != nil {
		http.Error(w, "something very wrong", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(b)
	w.WriteHeader(http.StatusOK)
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cant read data", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(b, &user)
	if err != nil {
		http.Error(w, "cant parse data", http.StatusBadRequest)
		return
	}

	err = h.auth.Register(user)
	if err != nil {
		http.Error(w, "cant register user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type loginData struct {
	Message      string `json:"message"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserRequest struct{}

type Handler struct {
	auth UserService
}

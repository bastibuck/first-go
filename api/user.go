package api

import (
	"encoding/json"
	"first-go/db"
	userTypes "first-go/types/user"
	"fmt"
	"net/http"
)

type UserHandler struct {
	userStore db.UserStore
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  *userTypes.User `json:"user"`
	Token string          `json:"token"`
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore,
	}
}

func (userHandler *UserHandler) RegisterUser(res http.ResponseWriter, req *http.Request) {
	var createUser userTypes.NewUserPayload

	ctx := req.Context()

	err := json.NewDecoder(req.Body).Decode(&createUser)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := userTypes.NewUser(createUser)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Users/Create", http.StatusInternalServerError)
		return
	}

	err = userHandler.userStore.Create(ctx, user)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Users/Create", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func (userHandler *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginParams LoginPayload

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&loginParams); err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := userHandler.userStore.GetByEmail(ctx, loginParams.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !userTypes.ValidatePassword(user.PasswordHash, loginParams.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := userTypes.CreateToken(*user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid credentials", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		User:  user,
		Token: token,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid credentials", http.StatusInternalServerError)
	}
}

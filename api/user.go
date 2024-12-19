package api

import (
	"encoding/json"
	"first-go/db"
	userTypes "first-go/types/user"
	"first-go/utils"
	"fmt"
	"net/http"
)

type UserHandler struct {
	userStore db.UserStore
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
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
	var createUser userTypes.RegisterPayload

	ctx := req.Context()

	err := json.NewDecoder(req.Body).Decode(&createUser)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	validate := utils.GetValidator()
	err = validate.Struct(createUser)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Malformed email or password", http.StatusBadRequest)
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

func (userHandler *UserHandler) LoginUser(res http.ResponseWriter, req *http.Request) {
	var loginParams LoginPayload

	ctx := req.Context()

	if err := json.NewDecoder(req.Body).Decode(&loginParams); err != nil {
		fmt.Println(err)
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	validate := utils.GetValidator()
	err := validate.Struct(loginParams)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := userHandler.userStore.GetByEmail(ctx, loginParams.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !userTypes.ValidatePassword(user.PasswordHash, loginParams.Password) {
		http.Error(res, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := userTypes.CreateToken(*user)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Invalid credentials", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		User:  user,
		Token: token,
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(response)

	if err != nil {
		fmt.Println(err)
		http.Error(res, "Invalid credentials", http.StatusInternalServerError)
	}
}

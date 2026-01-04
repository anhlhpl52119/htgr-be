package api

import (
	"encoding/json"
	"errors"
	"htrr-apis/internal/store"
	"htrr-apis/internal/utils"
	"log"
	"net/http"
	"regexp"
)

type UserHandler struct {
	logger    *log.Logger
	userStore store.UserStore
}

func NewUserHandler(logger *log.Logger, userStore store.UserStore) *UserHandler {
	return &UserHandler{
		logger:    logger,
		userStore: userStore,
	}
}

type registerUserRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("email is not valid")
	}

	if len(req.Phone) < 3 {
		return errors.New("phone is not valid")
	}

	return nil
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid payload"})
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: "",
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	err = h.userStore.CreateUser(user)

	if err != nil {
		h.logger.Printf("ERROR: registering user %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}

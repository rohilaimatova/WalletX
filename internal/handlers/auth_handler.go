package handlers

import (
	"WalletX/internal/service"
	"WalletX/models"
	"WalletX/pkg/errs"
	"WalletX/pkg/logger"
	"WalletX/pkg/respond"
	"WalletX/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Ping godoc
// @Summary      Health check
// @Description  Check if service is running
// @Tags         System
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /ping [get]
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("Ping endpoint called")
	respond.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "service is running",
	})
}

type UserHandler struct {
	Service        *service.UserService
	Redis          *redis.Client
	AccountService *service.AccountService
}

func NewUserHandler(s *service.UserService, accountSvc *service.AccountService, rdb *redis.Client) *UserHandler {
	return &UserHandler{
		Service:        s,
		AccountService: accountSvc,
		Redis:          rdb,
	}
}

// SignUp godoc
// @Summary Register user
// @Description Register user in two steps: 1) send SMS code 2) complete registration with code
// @Tags User
// @Accept json
// @Produce json
// @Param body body models.SignUpRequest true "Phone and optional SMS code"
// @Success 201 {object} models.RegisterResponse "Registration successful"
// @Success 201 {object} models.MessageResponse "SMS code sent"
// @Failure 400 {object} models.ErrorResponse "Invalid request or code"
// @Failure 409 {object} models.ErrorResponse "User already registered"
// @Router /api/users/signUp [post]
func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request", err)
		return
	}

	ctx := r.Context()

	if req.Code == "" {
		h.sendVerificationCode(w, ctx, req.Phone)
		return
	}

	h.registerUser(w, ctx, req.Phone, req.Code)
}

func (h *UserHandler) sendVerificationCode(w http.ResponseWriter, ctx context.Context, phone string) {
	user, _ := h.Service.GetByPhone(phone)
	if user.ID != 0 {
		logger.Info.Printf("Attempt to register already existing user: %s", phone)
		respond.JSON(w, http.StatusConflict, map[string]string{
			"message": "user already registered",
		})
		return
	}

	code := utils.GenerateCode()
	if err := h.Redis.Set(ctx, "verify:"+phone, code, 5*time.Minute).Err(); err != nil {
		logger.Error.Printf("Failed to save verification code for %s: %v", phone, err)
		respond.Error(w, http.StatusInternalServerError, "failed to save code", err)
		return
	}

	logger.Info.Printf("SMS verification code sent to %s: %s", phone, code)
	respond.JSON(w, http.StatusCreated, map[string]string{
		"message": "registration code sent",
	})
}

func (h *UserHandler) registerUser(w http.ResponseWriter, ctx context.Context, phone, code string) {
	savedCode, err := h.Redis.Get(ctx, "verify:"+phone).Result()
	if err != nil {
		logger.Warn.Printf("Verification code expired or not found for %s", phone)
		respond.Error(w, http.StatusBadRequest, "code expired or not found", errors.New("verification code expired or not found"))
		return
	}
	if savedCode != code {
		logger.Warn.Printf("Invalid verification code for %s", phone)
		respond.Error(w, http.StatusBadRequest, "invalid code", errors.New("invalid verification code"))
		return
	}

	user, err := h.Service.RegisterUser(phone)
	if err != nil {
		if errors.Is(err, errs.ErrUserExists) {
			logger.Info.Printf("User already exists: %s", phone)
			respond.JSON(w, http.StatusConflict, map[string]string{
				"message": "user already registered",
			})
			return
		}
		logger.Error.Printf("Failed to register user %s: %v", phone, err)
		respond.Error(w, http.StatusInternalServerError, "failed to create user", err)
		return
	}

	_, err = h.AccountService.CreateAccountForUser(user.ID)
	if err != nil {
		logger.Error.Printf("Failed to create account for userID=%d: %v", user.ID, err)
		respond.Error(w, http.StatusInternalServerError, "failed to create account", err)
		return
	}
	h.Redis.Del(ctx, "verify:"+phone)
	logger.Info.Printf("User registered successfully: %d", user.ID)
	respond.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "user registered",
		"user_id": user.ID,
	})
}

// SetPassword godoc
// @Summary      Set user password
// @Description  Sets password and returns JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.SetPasswordRequest true "Set password request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Router       /api/users/set-password [post]
func (h *UserHandler) SetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.SetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request", err)
		return
	}

	err := h.Service.SetPassword(req.UserID, req.Password)
	if err != nil {
		logger.Error.Printf("Failed to set password for user %d: %v", req.UserID, err)
		respond.HandleError(w, err)
		return
	}

	logger.Info.Printf("Password set successfully for user %d", req.UserID)

	token, err := utils.GenerateToken(req.UserID, "")
	if err != nil {
		logger.Error.Printf("Failed to generate token for user %d: %v", req.UserID, err)
		respond.Error(w, http.StatusInternalServerError, "failed to generate token", err)
		return
	}

	respond.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "password set successfully",
		"token":   token,
	})
}

// Login godoc
// @Summary      User login
// @Description  Login using phone and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Login request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /api/users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request", err)
		return
	}

	user, err := h.Service.Login(req.Phone, req.Password)
	if err != nil {
		logger.Warn.Printf("Failed login attempt for phone %s: %v", req.Phone, err)
		respond.HandleError(w, err)
		return
	}

	token, err := utils.GenerateToken(user.ID, "")
	if err != nil {
		logger.Error.Printf("Failed to generate token for user %d: %v", user.ID, err)
		respond.Error(w, http.StatusInternalServerError, "failed to generate token", err)
		return
	}

	logger.Info.Printf("User %d logged in successfully", user.ID)

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "logged in",
		"token":   token,
	})
}

// VerifyIdentity godoc
// @Summary      Verify user identity
// @Description  Verify user passport and personal data
// @Tags         User
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.VerifyIdentityRequest true "Identity data"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /api/users/verify [post]
func (h *UserHandler) VerifyIdentity(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyIdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request", err)
		return
	}

	err := h.Service.VerifyUser(req.UserID, req.FirstName, req.LastName, req.MiddleName, req.PassportNumber)
	if err != nil {
		logger.Error.Printf("Failed to verify identity for user %d: %v", req.UserID, err)
		respond.HandleError(w, err)
		return
	}

	logger.Info.Printf("User %d verified identity successfully", req.UserID)

	respond.JSON(w, http.StatusCreated, map[string]string{
		"message": "identity verified successfully",
	})
}

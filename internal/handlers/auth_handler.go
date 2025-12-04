package handlers

import (
	"WalletX/internal/service"
	"WalletX/models"
	"WalletX/pkg/errs"
	"WalletX/pkg/logger"
	"WalletX/pkg/utils"
	"WalletX/respond"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("Ping endpoint called")
	respond.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "service is running",
	})
}

type UserHandler struct {
	Service *service.UserService
	Redis   *redis.Client
}

func NewUserHandler(s *service.UserService, rdb *redis.Client) *UserHandler {
	return &UserHandler{
		Service: s,
		Redis:   rdb,
	}
}

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

// Вспомогательная функция для отправки кода
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

// Вспомогательная функция для регистрации пользователя
func (h *UserHandler) registerUser(w http.ResponseWriter, ctx context.Context, phone, code string) {
	savedCode, err := h.Redis.Get(ctx, "verify:"+phone).Result()
	if err != nil {
		logger.Warn.Printf("Verification code expired or not found for %s", phone)
		respond.Error(w, http.StatusBadRequest, "code expired or not found", nil)
		return
	}
	if savedCode != code {
		logger.Warn.Printf("Invalid verification code for %s", phone)
		respond.Error(w, http.StatusBadRequest, "invalid code", nil)
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

	h.Redis.Del(ctx, "verify:"+phone)
	logger.Info.Printf("User registered successfully: %d", user.ID)
	respond.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "user registered",
		"user_id": user.ID,
	})
}
func (h *UserHandler) SetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.SetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request", err)
		return
	}

	// 1️⃣ Устанавливаем пароль
	err := h.Service.SetPassword(req.UserID, req.Password)
	if err != nil {
		logger.Error.Printf("Failed to set password for user %d: %v", req.UserID, err)
		respond.HandleError(w, err)
		return
	}

	logger.Info.Printf("Password set successfully for user %d", req.UserID)

	// 2️⃣ Генерируем токен
	token, err := utils.GenerateToken(req.UserID, "")
	if err != nil {
		logger.Error.Printf("Failed to generate token for user %d: %v", req.UserID, err)
		respond.Error(w, http.StatusInternalServerError, "failed to generate token", err)
		return
	}

	// 3️⃣ Отправляем ответ
	respond.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "password set successfully",
		"token":   token,
	})
}
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

	// Генерируем токен
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

package handlers

import (
	"WalletX/internal/service"
	"WalletX/models"
	"WalletX/pkg/errs"
	"WalletX/pkg/utils"
	"WalletX/respond"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
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

	// ----------------------------------------------------
	// 1️⃣ Первый шаг — номер без кода → отправить SMS
	// ----------------------------------------------------
	if req.Code == "" {

		// Проверка — не зарегистрирован ли уже пользователь
		user, _ := h.Service.GetByPhone(req.Phone)
		if user.ID != 0 {
			respond.JSON(w, http.StatusConflict, map[string]string{
				"message": "user already registered",
			})
			return
		}

		// Генерируем код
		code := utils.GenerateCode()

		// сохраняем в redis на 5 минут
		err := h.Redis.Set(ctx, "verify:"+req.Phone, code, 5*time.Minute).Err()
		if err != nil {
			respond.Error(w, http.StatusInternalServerError, "failed to save code", err)
			return
		}

		fmt.Println("SMS to", req.Phone, "code:", code)

		respond.JSON(w, http.StatusCreated, map[string]string{
			"message": "registration code sent",
		})
		return
	}

	// ----------------------------------------------------
	// 2️⃣ Второй шаг — проверка кода
	// ----------------------------------------------------
	savedCode, err := h.Redis.Get(ctx, "verify:"+req.Phone).Result()
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "code expired or not found", nil)
		return
	}

	if savedCode != req.Code {
		respond.Error(w, http.StatusBadRequest, "invalid code", nil)
		return
	}

	// ----------------------------------------------------
	// 3️⃣ Создание пользователя
	// ----------------------------------------------------
	user, err := h.Service.RegisterUser(req.Phone)
	if err != nil {

		// Если пользователь уже существует → вернуть 409
		if errors.Is(err, errs.ErrUserExists) {
			respond.JSON(w, http.StatusConflict, map[string]string{
				"message": "user already registered",
			})
			return
		}

		respond.Error(w, http.StatusInternalServerError, "failed to create user", err)
		return
	}

	// Удаляем код после успешной регистрации
	h.Redis.Del(ctx, "verify:"+req.Phone)

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
		respond.HandleError(w, err)
		return
	}

	// 2️⃣ Генерируем токен
	token, err := utils.GenerateToken(req.UserID, "") // username пока можно ""
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to generate token", err)
		return
	}

	// 3️⃣ Отправляем ответ
	respond.JSON(w, http.StatusOK, map[string]interface{}{
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
		respond.HandleError(w, err)
		return
	}

	// генерируем токен
	token, err := utils.GenerateToken(user.ID, "")
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to generate token", err)
		return
	}

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
		respond.HandleError(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, map[string]string{
		"message": "identity verified successfully",
	})
}

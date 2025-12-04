package service

import (
	"WalletX/internal/repository"
	"WalletX/models"
	"WalletX/pkg/errs"
	"WalletX/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

// Глобальная переменная для проверки номера телефона
var phoneRegex = regexp.MustCompile(`^\+\d{12}$`)

type UserService struct {
	Repo repository.UserRepository
}

// Создаём сервис
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// Получение пользователя по телефону
func (s *UserService) GetByPhone(phone string) (models.User, error) {
	logger.Info.Printf("GetByPhone called for phone: %s", phone)
	return s.Repo.GetByPhone(phone)
}

// Регистрация пользователя по номеру телефона
func (s *UserService) RegisterUser(phone string) (models.User, error) {
	// Проверка номера телефона
	if !phoneRegex.MatchString(phone) {
		logger.Warn.Printf("RegisterUser: invalid phone format: %s", phone)
		return models.User{}, errs.ErrInvalidPhone
	}

	// Проверка на существование пользователя
	existing, _ := s.Repo.GetByPhone(phone)
	if existing.ID != 0 {
		logger.Warn.Printf("RegisterUser: phone already registered: %s", phone)
		return models.User{}, errs.ErrUserExists
	}

	// Создаём нового пользователя
	user := models.User{
		Phone:      phone,
		IsVerified: false,
	}

	created, err := s.Repo.CreateUser(user)
	if err != nil {
		logger.Error.Printf("RegisterUser: failed to create user: %v", err)
		return models.User{}, errs.ErrInternal
	}

	logger.Info.Printf("RegisterUser: success, userID=%d", created.ID)
	return created, nil
}

// Установка пароля (хеширование)
func (s *UserService) SetPassword(userID int, password string) error {
	if len(password) != 8 {
		logger.Warn.Printf("SetPassword: weak password for userID=%d", userID)
		return errs.ErrWeakPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error.Printf("SetPassword: failed to hash password: %v", err)
		return errs.ErrInternal
	}

	if err := s.Repo.UpdatePassword(userID, string(hash)); err != nil {
		logger.Error.Printf("SetPassword: failed to update DB: %v", err)
		return errs.ErrInternal
	}

	logger.Info.Printf("SetPassword: success for userID=%d", userID)
	return nil
}

// Вход пользователя
func (s *UserService) Login(phone, password string) (*models.User, error) {
	user, err := s.Repo.GetByPhone(phone)
	if err != nil {
		logger.Warn.Printf("Login failed: user not found for phone %s", phone)
		return nil, errs.ErrUserNotFound
	}

	// Проверяем заблокирован ли
	if user.DeviceID {
		logger.Warn.Printf("Login blocked: user %d is blocked", user.ID)
		return nil, errs.ErrUserBlocked
	}

	// Сравниваем пароль
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		// Прибавляем попытку
		if err := s.Repo.IncrementPasswordAttempts(user.ID); err != nil {
			return nil, err
		}
		user.PasswordAttempts++

		// Если 3 или больше попыток → блокируем
		if user.PasswordAttempts >= 3 {
			if err := s.Repo.BlockUser(user.ID); err != nil {
				return nil, err
			}
			logger.Warn.Printf("Login blocked: user %d reached max attempts", user.ID)
			return nil, errs.ErrUserBlocked
		}

		logger.Warn.Printf("Login failed: wrong password for user %d", user.ID)
		return nil, errs.ErrWrongPassword
	}

	// Пароль правильный → сбросить попытки
	s.Repo.ResetPasswordAttempts(user.ID)
	logger.Info.Printf("User %d logged in successfully", user.ID)
	return &user, nil
}

// Идентификация пользователя
func (s *UserService) VerifyUser(userID int, firstName, lastName, middleName, passport string) error {
	if firstName == "" || lastName == "" || middleName == "" || passport == "" {
		logger.Warn.Printf("VerifyUser: required fields missing, userID=%d", userID)
		return errs.ErrRequiredFields
	}

	if err := s.Repo.UpdateVerification(userID, firstName, lastName, middleName, passport); err != nil {
		logger.Error.Printf("VerifyUser: failed to update DB: %v", err)
		return errs.ErrInternal
	}

	logger.Info.Printf("VerifyUser: success for userID=%d", userID)
	return nil
}

package service

import (
	"WalletX/internal/repository"
	"WalletX/models"
	"WalletX/pkg/errs"
	"WalletX/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type UserService struct {
	Repo repository.UserRepository
}

// Создаём сервис
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}
func (s *UserService) GetByPhone(phone string) (models.User, error) {
	return s.Repo.GetByPhone(phone)
}

// Регистрация пользователя по номеру телефона
func (s *UserService) RegisterUser(phone string) (models.User, error) {
	// Проверка номера телефона
	re := regexp.MustCompile(`^\+\d{12}$`)
	if !re.MatchString(phone) {
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

	// Обновляем пароль в БД
	if err := s.Repo.UpdatePassword(userID, string(hash)); err != nil {
		logger.Error.Printf("SetPassword: failed to update DB: %v", err)
		return errs.ErrInternal
	}

	logger.Info.Printf("SetPassword: success for userID=%d", userID)
	return nil
}
func (s *UserService) Login(phone, password string) (*models.User, error) {
	user, err := s.Repo.GetByPhone(phone)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	// проверяем заблокирован ли
	if user.DeviceID {
		return nil, errs.ErrUserBlocked
	}

	// сравниваем пароль
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {

		// прибавляем попытку
		err := s.Repo.IncrementPasswordAttempts(user.ID)
		if err != nil {
			return nil, err
		}

		// обновляем локальное значение для проверки блокировки
		user.PasswordAttempts++

		// если 3 или больше попыток → блокируем
		if user.PasswordAttempts >= 3 {
			err := s.Repo.BlockUser(user.ID)
			if err != nil {
				return nil, err
			}
			return nil, errs.ErrUserBlocked
		}

		return nil, errs.ErrWrongPassword
	}

	// пароль правильный → сбросить попытки
	s.Repo.ResetPasswordAttempts(user.ID)

	return &user, nil
}

// Идентификация пользователя
func (s *UserService) VerifyUser(userID int, firstName, lastName, middleName, passport string) error {
	// Проверка что поля не пустые
	if firstName == "" || lastName == "" || middleName == "" || passport == "" {
		logger.Warn.Printf("VerifyUser: required fields missing, userID=%d", userID)
		return errs.ErrRequiredFields
	}

	// Обновляем данные пользователя
	err := s.Repo.UpdateVerification(userID, firstName, lastName, middleName, passport)
	if err != nil {
		logger.Error.Printf("VerifyUser: failed to update DB: %v", err)
		return errs.ErrInternal
	}

	logger.Info.Printf("VerifyUser: success for userID=%d", userID)
	return nil
}

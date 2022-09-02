package service

import (
	"crypto/sha1"
	"errors"

	"fmt"

	"time"

	"github.com/golang-jwt/jwt"
	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/internal/repository"
	"github.com/sirupsen/logrus"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"        // соль добавляемая к паролю пользователей
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH" // Набор случайных байт для подписи токена (ключ подписи) - так-же исп при расшифровке токена
	tokenTTL   = 12 * time.Hour                 // время жизни токена
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"user_id"`
}

// Авторизация и аутентификация
type AuthService struct {
	repo repository.Authorization
}

// NewAuthService конструктор
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Test() {
	logrus.Info("info in auth")
	logrus.WithFields(logrus.Fields{"tag": "a tag svc"}).Info("An info message")
	s.repo.TestDB()
}

//CreateUser создание пользователя
func (s *AuthService) CreateUser(user core.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

// Для генерации токена нужно получить пользователя из базы
// если пользователя нет, вернуть ошибку
// в токен записывается id пользователя
func (s *AuthService) GenerateToken(login, password string) (string, error) {
	user, err := s.repo.GetUser(login, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), // токен перестает быть валидным через
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

//Используется middleware
//Достаем Id пользователя из токена
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

//generatePasswordHash генерирует хеш, добавляем соль, перчим
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

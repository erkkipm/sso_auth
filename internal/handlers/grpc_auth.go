package handlers

import (
	"context"
	"github.com/erkkipm/sso_auth/internal/models"
	"github.com/erkkipm/sso_auth/internal/storage"
	"github.com/erkkipm/sso_auth/pkg/jwtutil"
	ssov1 "github.com/erkkipm/sso_proto/gen/go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

// ServerAPI ... реализует интерфейс pb.AuthServer
type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	Store  *storage.Storage
	JWTKey string
}

// NewServerAPI ...
func NewServerAPI(s *storage.Storage, jwtKey string) *ServerAPI {
	return &ServerAPI{
		Store:  s,
		JWTKey: jwtKey,
	}
}

// Register ... Регистрация нвоого пользователя
func (s *ServerAPI) Register(ctx context.Context, r *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	log.Printf("Register: входящий запрос: логин=%s телефон=%s Приложение=%s", r.Email, r.Phone, r.AppId)

	hash, _ := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	user := models.User{
		AppID:     r.AppId,
		Email:     r.Email,
		Password:  string(hash),
		Phone:     r.Phone,
		CreatedAt: time.Now(),
	}

	existing, err := s.Store.GetUserByEmailAndApp(ctx, user)
	if err != nil {
		log.Printf("Register: ошибка при обращении к базе: %v", err)
		return &ssov1.RegisterResponse{
			Success: false,
			Message: "Ошибка доступа! Сервис недоступен. Попробуйте снова чуток позже",
		}, err
	}
	if existing != nil {
		log.Printf("Register: ошибка при проверки пользователя на уникальноссть: %v", err)
		return &ssov1.RegisterResponse{
			Success: false,
			Message: "Пользователь с таким логином уже есть!",
		}, nil
	}

	if err := s.Store.CreateUser(ctx, user); err != nil {
		log.Printf("Register: ошибка при создании пользователя: %v", err)
		return &ssov1.RegisterResponse{
			Success: false,
			Message: "Ошибка подключения к базе: " + err.Error(),
		}, nil
	}

	log.Printf("Register: пользователь %s успешно создан", r.Email)

	return &ssov1.RegisterResponse{
		Success: true,
		Message: "Поздравляем! Вы зарегистрированы! В ближайшее время на e-mail придет письмо для подтверждения аккаунта",
	}, nil
}

// Login ...
func (s *ServerAPI) Login(ctx context.Context, r *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	user, err := s.Store.FindUser(ctx, r.AppId, r.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password)) != nil {
		return nil, err
	}
	token, _ := jwtutil.GenerateToken(user.Email, s.JWTKey, 72*time.Hour)
	return &ssov1.LoginResponse{Token: token}, nil
}

//// ChangePassword ...
//func (s *ServerAPI) ChangePassword(ctx context.Context, r *ssov1.) (*ssov1., error) {
//	user, err := s.Store.FindUser(ctx, r.AppId, r.Email)
//	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.OldPassword)) != nil {
//		return &ssov1.ChangePasswordResponse{Status: "error", Message: "invalid"}, nil
//	}
//	newHash, _ := bcrypt.GenerateFromPassword([]byte(r.NewPassword), bcrypt.DefaultCost)
//	if err := s.Store.UpdatePassword(ctx, r.AppId, r.Email, string(newHash)); err != nil {
//		return &ssov1.ChangePasswordResponse{Status: "error", Message: "ОШИБКА!"}, nil
//	}
//	return &ssov1.ChangePasswordResponse{Status: "ok", Message: "Пароль изменен"}, nil
//}

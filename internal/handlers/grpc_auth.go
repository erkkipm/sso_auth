package handlers

import (
	"context"
	ssoapb "github.com/erkkipm/sso_auth/gen/proto"
	"github.com/erkkipm/sso_auth/internal/models"
	"github.com/erkkipm/sso_auth/internal/storage"
	"github.com/erkkipm/sso_auth/pkg/jwtutil"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type AuthServer struct {
	ssoapb.UnimplementedAuthServiceServer
	Store  *storage.Storage
	JWTKey string
}

func NewAuthServer(s *storage.Storage, jwtKey string) *AuthServer {
	return &AuthServer{Store: s, JWTKey: jwtKey}
}

// Register ... Регистрация нвоого пользователя
func (a *AuthServer) Register(ctx context.Context, r *ssoapb.RegisterRequest) (*ssoapb.RegisterResponse, error) {
	log.Printf("Register: входящий запрос: логин=%s телефон=%s Приложение=%s", r.Email, r.Phone, r.AppId)

	hash, _ := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	user := models.User{
		AppID:     r.AppId,
		Email:     r.Email,
		Password:  string(hash),
		Phone:     r.Phone,
		CreatedAt: time.Now(),
	}

	existing, err := a.Store.GetUserByEmailAndApp(ctx, user)
	if err != nil {
		log.Printf("Register: ошибка при создании пользователя: %v", err)
		return &ssoapb.RegisterResponse{
			Status:  "error",
			Message: "Ошибка доступа! Сервис недоступен. Попробуйте снова чуток позже",
		}, err
	}
	if existing != nil {
		return &ssoapb.RegisterResponse{
			Status:  "error",
			Message: "Пользователь с таким логином уже есть!",
		}, nil
	}

	if err := a.Store.CreateUser(ctx, user); err != nil {
		log.Printf("Register: ошибка при создании пользователя: %v", err)
		return &ssoapb.RegisterResponse{
			Status:  "error",
			Message: "Ошибка подключения к базе: " + err.Error(),
		}, nil
	}

	log.Printf("Register: пользователь %s успешно создан", r.Email)

	return &ssoapb.RegisterResponse{
		Status:  "ok",
		Message: "Поздравляем! Пользователь зарегистрирован!",
	}, nil
}

func (a *AuthServer) Login(ctx context.Context, r *ssoapb.LoginRequest) (*ssoapb.LoginResponse, error) {
	user, err := a.Store.FindUser(ctx, r.AppId, r.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password)) != nil {
		return nil, err
	}
	token, _ := jwtutil.GenerateToken(user.Email, a.JWTKey, 72*time.Hour)
	return &ssoapb.LoginResponse{Token: token}, nil
}

func (a *AuthServer) ChangePassword(ctx context.Context, r *ssoapb.ChangePasswordRequest) (*ssoapb.ChangePasswordResponse, error) {
	user, err := a.Store.FindUser(ctx, r.AppId, r.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.OldPassword)) != nil {
		return &ssoapb.ChangePasswordResponse{Status: "error", Message: "invalid"}, nil
	}
	newHash, _ := bcrypt.GenerateFromPassword([]byte(r.NewPassword), bcrypt.DefaultCost)
	if err := a.Store.UpdatePassword(ctx, r.AppId, r.Email, string(newHash)); err != nil {
		return &ssoapb.ChangePasswordResponse{Status: "error", Message: "ОШИБКА!"}, nil
	}
	return &ssoapb.ChangePasswordResponse{Status: "ok", Message: "Пароль изменен"}, nil
}

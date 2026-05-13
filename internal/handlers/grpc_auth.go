package handlers

import (
	"context"
	"errors"
	"github.com/erkkipm/sso_auth/internal/models"
	"github.com/erkkipm/sso_auth/internal/storage"
	"github.com/erkkipm/sso_auth/pkg/jwtutil"
	ssoV1 "github.com/erkkipm/sso_proto/gen/go"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

// ServerAPI ... реализует интерфейс pb.AuthServer
type ServerAPI struct {
	ssoV1.UnimplementedAuthServer
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
func (s *ServerAPI) Register(ctx context.Context, r *ssoV1.RegisterRequest) (*ssoV1.RegisterResponse, error) {
	log.Printf("Register: входящий запрос: логин=%s телефон=%s Приложение=%s", r.Email, r.Phone, r.AppId)

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Register: ошибка хеширования пароля: %v", err)
		return &ssoV1.RegisterResponse{
			Success: false,
			Message: "Внутренняя ошибка сервера",
		}, status.Error(codes.Internal, "Ошибка хеширования пароля")
	}

	user := models.User{
		AppID:     r.AppId,
		Email:     r.Email,
		Password:  string(hash),
		Phone:     r.Phone,
		CreatedAt: time.Now(),
	}

	existing, err := s.Store.GetUserByEmailAndApp(ctx, user)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("Register: ошибка при обращении к базе: %v", err)
		return &ssoV1.RegisterResponse{
			Success: false,
			Message: "Ошибка доступа! Сервис недоступен. Попробуйте снова чуть позже",
		}, status.Error(codes.Internal, "Ошибка обращения к базе")
	}
	if existing != nil {
		log.Printf("Register: пользователь %s уже существует", r.Email)
		return &ssoV1.RegisterResponse{
			Success: false,
			Message: "Пользователь с таким email уже есть!",
		}, status.Error(codes.AlreadyExists, "Пользователь с таким email уже существует")
	}

	if err := s.Store.CreateUser(ctx, user); err != nil {
		log.Printf("Register: ошибка при создании пользователя: %v", err)
		return &ssoV1.RegisterResponse{
			Success: false,
			Message: "Ошибка подключения к базе: " + err.Error(),
		}, status.Error(codes.Internal, "Ошибка создания пользователя")
	}

	log.Printf("Register: пользователь %s успешно создан", r.Email)

	return &ssoV1.RegisterResponse{
		Success: true,
		Message: "Поздравляем! Вы зарегистрированы! В ближайшее время на e-mail придет письмо для подтверждения аккаунта",
	}, nil
}

// Login ...
func (s *ServerAPI) Login(ctx context.Context, req *ssoV1.LoginRequest) (*ssoV1.LoginResponse, error) {
	user, err := s.Store.FindUser(ctx, req.AppId, req.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("Login: ошибка поиска пользователя (%s): %v", req.Email, err)
		return nil, status.Error(codes.Internal, "Ошибка сервера, попробуйте позже")
	}
	if err == nil && user == nil {
		log.Printf("Login: пользователь не найден (%s)", req.Email)
		return nil, status.Error(codes.NotFound, "Пользователь с таким email не найден")
	}
	if user == nil {
		log.Printf("Login: пользователь не найден (%s)", req.Email)
		return nil, status.Error(codes.NotFound, "Пользователь с таким email не найден")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		log.Printf("Login: неверный пароль для %s.", req.Email)
		return nil, status.Error(codes.Unauthenticated, "Неверный пароль")
	}

	token, err := jwtutil.GenerateToken(user.ID.Hex(), user.Email, user.Phone, s.JWTKey, 1*time.Hour)
	if err != nil {
		log.Printf("Login: ошибка при генерации токена: %v", err)
		return nil, status.Error(codes.Internal, "Ошибка при генерации токена")
	}
	return &ssoV1.LoginResponse{
		Token: token,
	}, nil
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

package handlers

import (
	"context"
	"github.com/erkkipm/sso_auth/internal/models"
	"github.com/erkkipm/sso_auth/internal/storage"
	"github.com/erkkipm/sso_auth/pkg/jwtutil"
	"github.com/erkkipm/sso_auth/proto/proto"
	"golang.org/x/crypto/bcrypt"
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

func (a *AuthServer) Register(ctx context.Context, r *ssoapb.RegisterRequest) (*ssoapb.RegisterResponse, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	user := models.User{
		AppID:     r.AppId,
		Email:     r.Email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}
	if err := a.Store.CreateUser(ctx, user); err != nil {
		return &ssoapb.RegisterResponse{Status: "error", Message: "User exists"}, nil
	}
	return &ssoapb.RegisterResponse{Status: "ok", Message: "User registered"}, nil
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
		return &ssoapb.ChangePasswordResponse{Status: "error", Message: "fail"}, nil
	}
	return &ssoapb.ChangePasswordResponse{Status: "ok", Message: "updated"}, nil
}

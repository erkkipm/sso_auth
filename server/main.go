package main

import (
	"context"
	ssoapb "github.com/erkkipm/sso_auth/gen/proto"
	"github.com/erkkipm/sso_auth/internal/configs"
	"github.com/erkkipm/sso_auth/internal/handlers"
	"github.com/erkkipm/sso_auth/internal/logger"
	"github.com/erkkipm/sso_auth/internal/storage"
	"github.com/erkkipm/sso_auth/pkg/logger/sl"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// КОНТЕКСТ. Создание контекста с обработкой сигналов завершения
	_, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// КОНФИГУРАЦИЯ. Инициализация
	cfg := configs.GetConfig()

	// КОНТЕКСТ. ЗАПУСКАЕМ ЛОГГЕР
	log := logger.SetupLogger(cfg.Env, cfg.NameProject)
	log.Info("КОНФИГУРАЦИЯ. Успешно!", slog.String("project", cfg.NameProject))
	log.Info("ЛОГГЕР. Успешно!")
	log.Debug("ВКЛЮЧЕН РЕЖИМ ДЕБАГ!")

	// СОЗДАНИЕ ПРИЛОЖЕНИЯ
	log.Info("ПРИЛОЖЕНИЕ. Инициализация...", slog.String("env", cfg.Env))

	lis, err := net.Listen("tcp", ":"+cfg.HTTP.Port)
	if err != nil {
		log.Error("Ошибка порта:", sl.Err(err))
		os.Exit(1)
	}

	store, err := storage.NewStorage("mongodb://localhost:"+cfg.MongoDB.Port, cfg.MongoDB.Username, cfg.MongoDB.Collection.Users)
	if err != nil {
		log.Error("Mongo ошибка:", sl.Err(err))
		os.Exit(1)
	} else {
		log.Info("MongoDB. Успешно!", slog.String("mongodb://localhost:", cfg.MongoDB.Port+"/"+cfg.MongoDB.Username+"/"+cfg.MongoDB.Collection.Users))
	}

	s := grpc.NewServer()
	ssoapb.RegisterAuthServiceServer(s, handlers.NewAuthServer(store, cfg.JWTKey))

	log.Info("SSO_AUTH ПРИЛОЖЕНИЕ: Запущено!", slog.String("порт", cfg.HTTP.Port))

	if err := s.Serve(lis); err != nil {
		log.Error("ОШИБКА ПРИЛОЖЕНИЯ:", sl.Err(err))
		os.Exit(1)
	}
}

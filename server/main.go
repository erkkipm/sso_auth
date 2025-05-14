package main

import (
	"context"
	"github.com/erkkipm/sso_auth/internal/configs"
	"github.com/erkkipm/sso_auth/internal/handlers"
	"github.com/erkkipm/sso_auth/internal/logger"
	"github.com/erkkipm/sso_auth/internal/storage"
	"github.com/erkkipm/sso_auth/pkg/logger/sl"
	ssov1 "github.com/erkkipm/sso_proto/gen/go"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// КОНТЕКСТ. Создание контекста с обработкой сигналов завершения
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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

	// ПОДКЛЮЧЕНИЕ БД
	ctxMng, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	store, err := storage.NewStorage(ctxMng, cfg.MongoDB)
	if err != nil {
		log.Error("БД. Ошибка подключения к MongoDB:", sl.Err(err))
		os.Exit(1)
	} else {
		log.Info("БД. Успешно подключена!", slog.String("mongodb://localhost:", cfg.MongoDB.Port+"/"+cfg.MongoDB.Username+"/"+cfg.MongoDB.Collection.Users))
	}

	// Создаем TCP-листенер на порту 50055 для gRPC
	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Error("Ошибка подключения порта:", sl.Err(err))
		os.Exit(1)
	}
	// Создаем gRPC-сервер
	gRPCServer := grpc.NewServer()

	// Регистрируем наш сервис Users на gRPC-сервере
	ssov1.RegisterAuthServer(gRPCServer, handlers.NewServerAPI(store, cfg.JWTKey))
	log.Info("ПРИЛОЖЕНИЕ: Запущено!", slog.String("порт", cfg.GRPC.Port))
	// Запускаем сервер (блокирующий вызов)
	if err := gRPCServer.Serve(lis); err != nil {
		log.Error("ОШИБКА ПРИЛОЖЕНИЯ:", sl.Err(err))
		os.Exit(1)
	}
	// =========================================================================

}

package main

import (
	"context"
	ssoapb "github.com/erkkipm/sso_auth/gen/proto"
	"github.com/erkkipm/sso_auth/internal/configs"
	"github.com/erkkipm/sso_auth/internal/handlers"
	"github.com/erkkipm/sso_auth/internal/storage"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/signal"
	"syscall"
)

func main() {

	// КОНТЕКСТ. Создание контекста с обработкой сигналов завершения
	_, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// КОНФИГУРАЦИЯ. Инициализация
	cfg := configs.GetConfig()

	lis, err := net.Listen("tcp", ":"+cfg.HTTP.Port)
	if err != nil {
		log.Fatal("Ошибка порта:", err)
	}

	store, err := storage.NewStorage("mongodb://localhost:"+cfg.MongoDB.Port, cfg.MongoDB.Username, cfg.MongoDB.Collection.Users)
	if err != nil {
		log.Fatal("Mongo ошибка:", err)
	}

	s := grpc.NewServer()
	ssoapb.RegisterAuthServiceServer(s, handlers.NewAuthServer(store, cfg.JWTKey))

	log.Println("sso_auth gRPC сервер запущен на порту:", cfg.HTTP.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

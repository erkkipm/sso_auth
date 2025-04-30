package main

import (
	"github.com/erkkipm/sso_auth/internal/handlers"
	"github.com/erkkipm/sso_auth/internal/storage"
	"github.com/erkkipm/sso_auth/proto/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Ошибка порта:", err)
	}

	store, err := storage.NewStorage("mongodb://localhost:38838", "authsofit", "users")
	if err != nil {
		log.Fatal("Mongo ошибка:", err)
	}

	s := grpc.NewServer()
	ssoapb.RegisterAuthServiceServer(s, handlers.NewAuthServer(store, "ErkkiSofit9944"))

	log.Println("sso_auth gRPC сервер запущен на :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

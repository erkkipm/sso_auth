# Makefile для проекта sso_auth

# Название proto-файла
PROTO_FILE=proto/auth.proto

# Префикс go-пакета
GO_PACKAGE=github.com/erkkipm/sso_auth

# Папка для сгенерированных файлов
GEN_DIR=proto

all: generate run

# Генерация protobuf
.PHONY: generate
generate:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       --go_out=$(GEN_DIR) --go-grpc_out=$(GEN_DIR) \
	       $(PROTO_FILE)

# Запуск сервера
.PHONY: run
run:
	go run ./server/main.go

# Установка зависимостей
.PHONY: deps
deps:
	go mod tidy

# Очистка (по необходимости)
.PHONY: clean
clean:
	rm -f $(GEN_DIR)/*.pb.go

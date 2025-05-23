# Переменные окружения
NAME_APP = sso_auth
APP_MAIN_GO = ./server/main.go
#CONFIG_APP = ./configs/config_prod.yaml
#NAME_DB = sso_auth
#DB_USER = sso_auth
#DB_PASS = ErkkiSofit9944
DB_PORT = 38841
DB_PATH = /Users/erkkipuolakainen/data/db/
#GEN_DIR = ./gen/
#PROTO_DIR = ./proto/


# Установка зависимостей, сборка, запуск Mongo и создание пользователя
.PHONY: all
all: tidy mongo-start run
#all: mongo-stop clean build mongo-start mongo-create-user run
#all: mongo-stop tidy clean build mongo-start mongo-create-user gen run

# Установка зависимостей
.PHONY: tidy
tidy:
	@echo "======= Установка зависимостей ========"
	@go mod tidy && echo " ✅  Зависимости установлены!" || echo " ❌  Зависимости не установлены!"

# Очистка
#.PHONY: clean
#clean:
#	@echo "======= Удаление старых файлов ========"
#	@rm -rf ./build/ && echo " ✅  Старые файлы удалены!"

# Сборка
#.PHONY: build
#build:
#	@echo "======= Сборка приложения ========"
#	@mkdir -p ./build/
#	@go build -o ./build/sso_auth ./server/main.go && echo " ✅  Сборка успешна!" || echo " ❌  Ошибка сборки!"

# Остановка MongoDB
#.PHONY: mongo-stop
#mongo-stop:
#	@echo "======= Остановка MongoDB... ========"
#	@brew services stop mongodb/brew/mongodb-community && echo " ✅  MongoDB остановлен" || echo " ❌  MongoDB не найден."
#	@sudo mongosh --port $(DB_PORT) --eval "db.shutdownServer()" && \
#		echo " ❌  MongoDB на порту $(DB_PORT) пришлось остановить принудительно!" || \
#		echo " ✅  MongoDB не запущен!"

# Запуск MongoDB
.PHONY: mongo-start
mongo-start:
	@echo "======= Запуск MongoDB ========"
	@if nc -z 127.0.0.1 $(DB_PORT); then \
		echo "✅  MongoDB уже запущен на порту $(DB_PORT)."; \
	else \
		echo "✅  MongoDB не запущен. Пытаемся стартовать..."; \
		sudo mongod \
			--dbpath=$(DB_PATH) \
			--port=$(DB_PORT) \
			--logpath=$(DB_PATH)/mongod.log \
			--logappend \
			--fork \
			--bind_ip 127.0.0.1 || echo "❌ Не удалось запустить MongoDB"; \
		sleep 2; \
		if nc -z 127.0.0.1 $(DB_PORT); then echo "✅  MongoDB успешно запущен."; \
		else echo "❌  MongoDB не слушает порт. Последние строки лога:"; tail -n 20 $(DB_PATH)/mongod.log; false; fi; \
	fi

# Создание пользователя MongoDB
#.PHONY: mongo-create-user
#mongo-create-user:
#	@echo "======= Создание пользователя MongoDB ========"
#	@mongosh --port $(DB_PORT) --eval 'db.getSiblingDB("admin").createUser({user: "$(DB_USER)", pwd: "$(DB_PASS)", roles: [{role: "readWrite", db: "$(NAME_DB)"}]})' \
#	&& echo "✅ Пользователь $(DB_USER) создан в admin" \
#	|| echo "✅ Возможно, пользователь уже существует в admin."
#	@mongosh --port $(DB_PORT) --eval 'db.getSiblingDB("$(NAME_DB)").createUser({user: "$(DB_USER)", pwd: "$(DB_PASS)", roles: [{role: "readWrite", db: "$(NAME_DB)"}]})' \
#	&& echo "✅ Пользователь $(DB_USER) создан в базе $(NAME_DB)" \
#	|| echo "✅ Возможно, пользователь уже существует в $(NAME_DB)."
#
#	@echo "======= Проверка пользователя в admin ========"
#	@mongosh admin --port $(DB_PORT) --eval 'u = db.getUser("$(DB_USER)"); if (u) { printjson(u) } else { print("❌ Пользователь не найден в admin") }' \
#	|| echo "⚠️ Не удалось подключиться к MongoDB admin"
#
#	@echo "======= Проверка пользователя в базе $(NAME_DB) ========"
#	@mongosh $(NAME_DB) --port $(DB_PORT) --eval 'u = db.getUser("$(DB_USER)"); if (u) { printjson(u) } else { print("❌ Пользователь не найден в $(NAME_DB)") }' \
#	|| echo "⚠️ Не удалось подключиться к MongoDB $(NAME_DB)"


.PHONY: run
run:
	@echo "======= Запуск приложения ========"
	@go run $(APP_MAIN_GO)  && echo " ✅  Приложение запущено!" || echo " ❌  Приложение не запущено!"


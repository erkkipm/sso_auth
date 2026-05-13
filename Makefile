APP_BIN = ./build/app
APP_MAIN = ./server/main.go
APP_CONFIG = ./configs/config_prod.yaml

DB_PORT = 38841
HOME := $(shell echo $$HOME)
DB_USER = sso_auth
DB_PASS = ErkkiSofit9944
DB_PATH = $(HOME)/data/db
DB_LOG = ./log/mongo.log

USE_SUDO = sudo

.PHONY: all start install tidy build run check-ffprobe \
        mongodb-kill mongodb-clean mongodb-init-user mongodb-start mongodb

#  Запуск при запущенном монгодиби
start: run-app

# Полная сборка
run: build run

# ЖЕСТКИЙ РЕСТАРТ
install: tidy check-ffprobe mongodb

# Проверка ffprobe и ffmpeg
check-ffprobe:
	@echo "== Проверка ffprobe =="
	@if ! command -v ffprobe >/dev/null 2>&1; then \
		echo "⚠️  ffprobe не найден. Устанавливаем ffmpeg..."; \
		brew install ffmpeg || { echo "❌ Не удалось установить ffmpeg"; exit 1; }; \
	else \
		echo "✅ ffprobe уже установлен"; \
	fi
	@ffprobe -v error \
		-show_entries format=format_name \
		-of default=noprint_wrappers=1:nokey=1 \
		./test.mov || echo "⚠️  Файл ./test.mov не прошёл проверку"

# Установка зависимостей Go
tidy:
	@echo "======= Установка зависимостей ========"
	@go mod tidy && echo " ✅  Зависимости установлены!" || echo " ❌  Ошибка установки зависимостей"

# Сборка приложения
build:
	@echo "======= Сборка приложения ========"
	@mkdir -p build
	@mkdir -p ./log
	@go build -o $(APP_BIN) $(APP_MAIN) && echo " ✅  Сборка успешна!" || { echo " ❌  Ошибка сборки"; exit 1; }

# Запуск приложения
run:
	@echo "======= Запуск приложения ========"
	@env CONFIG-PATH=$(APP_CONFIG) $(APP_BIN)


run-app:
	@echo "🧪 Проверка доступности MongoDB на порту $(DB_PORT)..."
	@if nc -z 127.0.0.1 $(DB_PORT); then \
		echo "✅ MongoDB уже запущен."; \
	else \
		echo "⚠️  MongoDB не запущен. Запускаем MongoDB..."; \
		make mongodb-start || { echo '❌ Не удалось запустить MongoDB'; exit 1; }; \
	fi
	@echo "🧹 Удаляем старую сборку (если есть)..."
	@rm -f $(APP_BIN)
	@echo "🔨 Собираем приложение..."
	@make build
	@echo "🚀 Запускаем приложение..."
	@env CONFIG-PATH=$(APP_CONFIG) $(APP_BIN)


# MongoDB: полный цикл инициализации и запуска

# Убить mongod на нужном порту
mongodb-kill:
	@echo "⛔️ Завершаем все mongod на порту $(DB_PORT)..."
	@sudo rm -f /tmp/mongodb-$(DB_PORT).sock
	@lsof -ti tcp:$(DB_PORT) | xargs -r kill
	@ps aux | grep mongod || true

# Почистить lock-файл
mongodb-clean:
	@if [ -f $(DB_PATH)/$(DB_USER)/mongod.lock ]; then \
		rm -f $(DB_PATH)/$(DB_USER)/mongod.lock; \
	fi
	@if [ -f $(DB_PATH)/$(DB_USER)/storage.bson ]; then \
		rm -f $(DB_PATH)/$(DB_USER)/storage.bson; \
	fi
	@if [ -d $(DB_PATH)/$(DB_USER) ]; then \
		sudo chown -R $(USER):staff $(DB_PATH)/$(DB_USER); \
	fi

# Инициализация пользователя (запуск без авторизации, создание юзера, остановка)
mongodb-init-user: mongodb-kill mongodb-clean
	@echo "➡️  Запускаем mongod без авторизации для инициализации пользователя..."
	@mkdir -p $(DB_PATH)/$(DB_USER)/ && echo "Папка для хранилища MongoDB создана" || echo "ОШИБКА! Папка не создана"
	@rm -rf ./log && mkdir -p ./log
	@rm -f $(DB_PATH)/$(DB_USER)/mongod.lock
	@rm -f ./log/mongo.log
	@mongod --port $(DB_PORT) --dbpath $(DB_PATH)/$(DB_USER)/ --bind_ip 127.0.0.1 --noauth --fork --logpath $(DB_LOG) --logappend
	@sleep 2
	@echo "➡️  Проверяем наличие пользователя $(DB_USER)..."
	@if mongosh admin --port $(DB_PORT) --quiet --eval '!!db.getUser("$(DB_USER)")' 2>/dev/null | grep -q true; then \
		echo "⚠️  Пользователь \033[1m$(DB_USER)\033[0m уже существует в MongoDB на порту $(DB_PORT)."; \
		echo "❗️ ВСЕ данные в $(DB_PATH)/$(DB_USER)/ будут безвозвратно удалены."; \
		read -p "❓ Вы уверены, что хотите продолжить? (yes/[no]): " confirm; \
		if [ "$$confirm" = "yes" ]; then \
			echo "🧨 Удаляем старую базу данных..."; \
			$(USE_SUDO) pkill -f "mongod.*$(DB_PORT)" || true; \
			sleep 2; \
			$(USE_SUDO) rm -rf $(DB_PATH)/$(DB_USER)/; \
			mkdir -p $(DB_PATH)/$(DB_USER)/; \
			mongod --port $(DB_PORT) --dbpath $(DB_PATH)/$(DB_USER)/ --bind_ip 127.0.0.1 --noauth --fork --logpath $(DB_LOG) --logappend; \
			sleep 2; \
		else \
			echo "🚫 Отмена. База данных и пользователь сохранены."; \
			pkill -f "mongod.*$(DB_PORT)" || true; \
			exit 1; \
		fi \
	fi
	@mongosh admin --port $(DB_PORT) --eval 'db.createUser({user: "$(DB_USER)", pwd: "$(DB_PASS)", roles: [{ role: "readWrite", db: "$(DB_USER)" }, { role: "readWrite", db: "admin" }]})'
	@mongosh admin --port $(DB_PORT) --eval 'db.getUsers()'
	@pkill -f "mongod.*$(DB_PORT)" || true
	@sleep 2
	@echo "✅ Пользователь создан!"


# Запуск mongod с авторизацией
mongodb-start: mongodb-clean
	@mongod --port $(DB_PORT) --dbpath $(DB_PATH)/$(DB_USER)/ \
		--bind_ip 127.0.0.1 --auth --fork --logpath $(DB_LOG) --logappend || { \
		echo "❌ Ошибка запуска MongoDB. Вывод лога:"; \
		cat $(DB_LOG); exit 1; }
		@sleep 2
	@echo "✅ MongoDB запущен с авторизацией на порту $(DB_PORT)"

# Основная задача для базы: kill, clean, create user, старт
mongodb: mongodb-init-user mongodb-start



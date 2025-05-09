DB_PORT = 38839
SERVER_IP = 127.0.0.1
SERVER_PORT = 50051
DB_USER = sso_auth
APP_NAME = sso_auth
DB_PASS = ErkkiSofit9944
APP_BIN = ./build/
APP_MAIN = ./server/main.go
CONFIG_APP = ./configs/config_prod.yaml
ROOT_DIR = /Users/erkkipuolakainen/go/src/github.com/erkkipm/sso_auth/
PATH_SERVICE = /etc/systemd/system/
DB_PATH = /Users/erkkipuolakainen/data/db/
PROTO_DIR = ./proto/
GEN_DIR = ./gen/
GO_PACKAGE = github.com/erkkipm/sso_auth

.PHONY: all
all: stop deps clean gen build start get-users run

.PHONY: run
run:
	echo "======= Запуск приложения ========"
	go run $(APP_MAIN) && echo " ✅  Приложение запущено!" || echo " ❌  Приложение не запущено!"

.PHONY: gen
gen:
	echo "======= Генерация кода ========"
	mkdir -p $(GEN_DIR)
	protoc --go_out=$(GEN_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_DIR)/*.proto && echo " ✅  Код сгенерирован!" || echo " ❌  Код не сгенерирован!"

.PHONY: build
build:
	echo "======= Сборка приложения ========"
	mkdir -p $(APP_BIN)
	go build -o $(APP_BIN)/$(APP_NAME) $(APP_MAIN) && echo " ✅  Приложение собрано!" || echo " ❌  Приложение не собрано!"


.PHONY: clean
clean:
	echo "======= Удаление старых файлов ========"
	rm -rf $(APP_BIN) && echo " ✅  "$(APP_BIN)" удалено!" || echo " ❌  "$(APP_BIN)" не удалено!"
	rm -f $(GEN_DIR)/*.pb.go && echo " ✅  "$(GEN_DIR)" удалено!" || echo " ❌  "$(GEN_DIR)" не удалено!"

.PHONY: stop start create-user get-users restart
stop:
	echo "======= Остановка MongoDB... ========"
	brew services stop mongodb/brew/mongodb-community && echo " ✅  MongoDB остановлен" || echo " ❌  MongoDB не найден."
	sudo mongosh --port $(DB_PORT) --eval "db.shutdownServer()" && echo " ❌  MongoDB на порту "$(DB_PORT)" пришлось остановить принудительно!" || echo " ✅  MongoDB не запущен!"
#	sudo pkill mongod || echo "ОШИБКА!"
start:
	echo "======= Запуск MongoDB... ========"
	touch /var/log/mongodb.log || true
	mkdir -p $(DB_PATH)$(DB_USER)/ && echo " ✅  Папка для хранилища MongoDB создана" || echo " ❌  Папка для хранилища MongoDB не создана"
	sudo mongod --port $(DB_PORT) --dbpath=$(DB_PATH)$(DB_USER)/ --noauth --fork --logpath /var/log/mongodb.log --logappend && echo " ✅  MongoDB на порту "$(DB_PORT)" ЗАПУЩЕН"
create-user:
	echo "======= Создание пользователей MongoDB... ======="
	sudo mongosh admin --port $(DB_PORT) --eval 'if (!db.getUser("$(DB_USER)")) { db.createUser({ user: "$(DB_USER)", pwd: "$(DB_PASS)", roles: [{ role: "readWrite", db: "admin" }] }) }'
	sudo mongosh $(DB_USER) --port $(DB_PORT) --eval 'if (!db.getUser("$(DB_USER)")) { db.createUser({ user: "$(DB_USER)", pwd: "$(DB_PASS)", roles: [{ role: "readWrite", db: "$(DB_USER)" }] }) }'  && echo " ✅  Пользователи созданы!" || echo " ❌  Пользователи не созданы!"
get-users:
	echo "======= Получение пользователей MongoDB... ========"
	sudo mongosh admin --port $(DB_PORT) --eval 'db.getUsers()'
	sudo mongosh $(DB_USER) --port $(DB_PORT) --eval 'db.getUsers()'
restart:
	echo "======= Перезапуск MongoDB... ========"
	sudo mongosh --port $(DB_PORT) --eval "db.shutdownServer()" && echo " ❌  MongoDB на порту "$(DB_PORT)" пришлось остановить принудительно!" || echo " ✅  MongoDB не запущен!"
	sudo mongod --port $(DB_PORT) --dbpath=$(DB_PATH)$(DB_USER)/ --bind_ip 0.0.0.0 --auth --fork --logpath /var/log/mongodb.log --logappend  && echo " ✅  MongoDB перезапущен!" || echo " ❌  MongoDB не перезапущен!"

.PHONY: deps
deps:
	echo "======= Установка зависимостей ========"
	go mod tidy && echo " ✅  Зависимости установлены!" || echo " ❌  Зависимости не установлены!"

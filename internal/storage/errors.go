package storage

import "errors"

var (
	// ErrRecordNotFound ... Если запись в БД не найдена
	ErrRecordNotFound = errors.New("Запись не найдена")
	ErrNoRecord       = errors.New("Нет ни одной записи")

	// ErrUserExists ... ПОЛЬЗОВАТЕЛЬ. Если уже существует в базе данных
	ErrUserExists = errors.New("Пользователь уже существует")
	// ErrUserNotFound ... ПОЛЬЗОВАТЕЛЬ. Если пользователь не найден
	ErrUserNotFound = errors.New("Пользователь не найден")
	// ErrAppNotFound ... ПРИЛОЖЕНИЕ. Если не найдено
	ErrAppNotFound = errors.New("Приложение не найдено")

	ErrInvalidCredentials = errors.New("Неправильные учетные данные")
)

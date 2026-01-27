package linkusecase

import "errors"

var (
	/*Не найден*/
	ErrNotFound          = errors.New("link not found")
	/*Конфликт*/
	ErrShortNameConflict = errors.New("short_name already exists")
	/*Невалидный ввод*/
	ErrInvalidInput      = errors.New("invalid input")
)

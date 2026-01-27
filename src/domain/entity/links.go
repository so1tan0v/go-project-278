package entity

import "time"

/*Entity для ссылок*/
type Link struct {
	ID          int64     /*Идентифиактор записи*/
	OriginalURL string    /*Исходный URL*/
	ShortName   string    /*Короткий URL*/
	CreatedAt   time.Time /*Дата создания*/
}

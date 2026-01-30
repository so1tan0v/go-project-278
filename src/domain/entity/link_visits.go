package entity

import "time"

/*Entity для посещений ссылок*/
type LinkVisit struct {
	ID        int64     /*Идентифиактор записи*/
	LinkID    int64     /*Идентификатор ссылки*/
	IP        string    /*IP клиента*/
	UserAgent string    /*User-Agent*/
	Referer   string    /*Referer*/
	Status    int       /*HTTP статус редиректа*/
	CreatedAt time.Time /*Дата создания*/
}


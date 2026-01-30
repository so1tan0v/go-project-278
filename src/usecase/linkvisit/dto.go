package linkvisitusecase

import "time"

/*DTO для посещения ссылки*/
type LinkVisitDTO struct {
	ID        int64     `json:"id"`
	LinkID    int64     `json:"link_id"`
	CreatedAt time.Time `json:"created_at"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
	Status    int       `json:"status"`
}


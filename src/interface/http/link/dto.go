package link

/*DTO для ответа API.*/
type LinkResponse struct {
	ID          int64  `json:"id"`           /*Идентификатор ссылки*/
	OriginalURL string `json:"original_url"` /*Исходная ссылка*/
	ShortName   string `json:"short_name"`   /*Короткая ссылка*/
	ShortURL    string `json:"short_url"`    /*Короткая ссылка*/
}

/*DTO для создания ссылки.*/
type CreateLinkRequest struct {
	OriginalURL string `json:"original_url"` /*Исходная ссылка*/
	ShortName   string `json:"short_name"`   /*Короткая ссылка*/
}

/*DTO для обновления ссылки.*/
type UpdateLinkRequest struct {
	OriginalURL string `json:"original_url"` /*Исходная ссылка*/
	ShortName   string `json:"short_name"`   /*Короткая ссылка*/
}

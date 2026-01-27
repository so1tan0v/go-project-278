package linkusecase

/*DTO для работы с ссылками*/
type LinkDTO struct {
	ID          int64  /*Идентификатор ссылки*/
	OriginalURL string /*Исходная ссылка*/
	ShortName   string /*Короткая ссылка*/
	ShortURL    string /*Короткая ссылка*/
}

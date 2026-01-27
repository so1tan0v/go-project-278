package linkusecase

// LinkDTO — данные о ссылке, возвращаемые из UseCase.
type LinkDTO struct {
	ID          int64
	OriginalURL string
	ShortName   string
	ShortURL    string
}

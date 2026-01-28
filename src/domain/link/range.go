package link

/*Параметры для парсинга range из заголовка HTTP*/
type Range struct {
	Start int `json:"start"` /*Начало*/
	End   int `json:"end"`   /*Конец*/
}

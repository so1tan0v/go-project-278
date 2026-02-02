package configDomain

/*Конфигурация приложения*/
type AppConfig struct {
	Development bool   /*Разработка*/
	Port        int    /*Порт*/
	Host        string /*Хост*/
	LoggingIO   bool   /*Логирование*/
	BaseURL     string /*Базовый URL*/
	AllowedOrigins []string /*Разрешенные Origin*/
}
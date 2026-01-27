package configDomain

/*Конфигурация приложения*/
type Config struct {
	App      AppConfig      /*Конфигурация приложения*/
	Database DatabaseConfig /*Конфигурация базы данных*/
}

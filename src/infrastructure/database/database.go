package database

import (
	"database/sql"
	"fmt"

	"link-service/src/infrastructure/database/postgres"

	configDomain "link-service/src/domain/config"
)

/*Экземпляр базы данных*/
type DatabaseImpl struct {
	instance  *sql.DB
	config    configDomain.DatabaseConfig
	loggingIO bool
}

/*Метод создания нового экземпляра базы данных*/
func NewDatabaseImpl(config configDomain.DatabaseConfig, loggingIO bool) *DatabaseImpl {
	return &DatabaseImpl{
		config:    config,
		loggingIO: loggingIO,
	}
}

/*Метод подключения к базе данных*/
func (d *DatabaseImpl) Connect(config configDomain.DatabaseConfig) error {
	db, err := postgres.Open(config.URL)
	if err != nil {
		return err
	}

	d.instance = db

	if d.loggingIO {
		fmt.Println("Database connected successfully")
	}

	return nil
}

/*Метод получения экземпляра базы данных*/
func (d *DatabaseImpl) GetInstance() any { // *sql.DB
	return d.instance
}

/*Метод отключения от базы данных*/
func (d *DatabaseImpl) Disconnect() error {
	if d.instance == nil {
		return nil
	}

	err := d.instance.Close()
	if d.loggingIO {
		fmt.Println("Database disconnected successfully")
	}

	return err
}

/*Метод проверки соединения с базой данных*/
func (d *DatabaseImpl) Ping() error {
	if d.instance == nil {
		return fmt.Errorf("db not connected")
	}

	return d.instance.Ping()
}

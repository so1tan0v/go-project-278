package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"link-service/src/config"
	configDomain "link-service/src/domain/config"
	database "link-service/src/infrastructure/database"
	postgreslinkrepo "link-service/src/infrastructure/repository/postgres"
	httpinterface "link-service/src/interface/http"
	linkusecase "link-service/src/usecase/link"
	linkvisitusecase "link-service/src/usecase/linkvisit"
)

func main() {
	var httpServer *gin.Engine
	var cnf *configDomain.Config
	var db *database.DatabaseImpl

	var err error

	httpServer = gin.Default()

	cnf, err = config.Init(os.Getenv("ENV_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	httpServer.Use(cors.New(cors.Config{
		AllowOrigins: cnf.App.AllowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{
			"Content-Range",
		},
		MaxAge: 12 * time.Hour,
	}))

	db = database.NewDatabaseImpl(cnf.Database, cnf.App.LoggingIO)
	if err := db.Connect(cnf.Database); err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer func() { _ = db.Disconnect() }()

	sqlDB, ok := db.GetInstance().(*sql.DB)
	if !ok || sqlDB == nil {
		log.Fatal("invalid database instance")
	}

	linkRepo := postgreslinkrepo.New(sqlDB)
	linkService := linkusecase.NewService(linkRepo, cnf.App.BaseURL)

	linkVisitRepo := postgreslinkrepo.NewLinkVisitRepository(sqlDB)
	linkVisitService := linkvisitusecase.NewService(linkVisitRepo)

	httpServer.Use(gin.Recovery())
	httpinterface.InitRoutes(httpServer, httpinterface.Deps{
		Link:      linkService,
		LinkVisit: linkVisitService,
	})

	if err := httpServer.Run(fmt.Sprintf("%s:%d", cnf.App.Host, cnf.App.Port)); err != nil {
		log.Fatal(err)
	}
}

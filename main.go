package main

import (
	"encoding/json"
	"fmt"
	"link-service/src/config"
	httpInterface "link-service/src/interfaces/http"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	httpServer := gin.Default()

	cnf, err := config.Init(os.Getenv("ENV_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	if cnf.Development {
		httpServer.Use(gin.Recovery())

		jsonData, _ := json.Marshal(cnf)

		fmt.Println(string(jsonData))
	}

	httpInterface.InitRoutes(httpServer)

	if err := httpServer.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

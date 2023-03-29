package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"testApplication/graph"
	"testApplication/handlers"
	"testApplication/interfaces"
	"testApplication/repositories/mongodb"
	"testApplication/repositories/postgres"
	"testApplication/utils"
)

func main() {

	utils.LoadConf()

	usingDatabase := utils.Conf.Get("usingDatabase")

	var repo interfaces.ClientRepo

	switch usingDatabase {
	case "postgres":
		repo = postgres.InitConnection()
	case "mongo":
		repo = mongodb.InitConnection()
	default:
		log.Fatal("Wrong value for usingDatabase parameter, check config")
	}

	handler, _ := handlers.NewClientHandler(repo)
	router := gin.Default()
	router.GET("/clients", handler.GetClients)
	router.GET("/clients/:id", handler.GetClientById)
	router.POST("/clients", handler.CreateClient)
	router.PATCH("/clients", handler.UpdateClient)
	router.DELETE("/clients/:id", handler.DeleteClient)

	newGraph, err := graph.NewGraph(repo)
	if err != nil {
		return
	}

	router.POST("/graph", newGraph.GraphqlHandler)

	router.Run("127.0.0.1:8080")
}

package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"testApplication/graph"
	"testApplication/handlers"
	"testApplication/interfaces"
	"testApplication/middleware"
	"testApplication/redis"
	"testApplication/repositories/mongodb"
	"testApplication/repositories/postgres"
	"testApplication/utils"
)

func main() {

	utils.LoadConf()

	redisToken, err := redis.NewConn()
	if err != nil {
		panic(err)
	}

	usingDatabase := utils.Conf.Get("usingDatabase")

	var repoClient interfaces.ClientRepo
	repoUsers := postgres.InitConnection()

	switch usingDatabase {
	case "postgres":
		repoClient = postgres.InitConnection()
	case "mongo":
		repoClient = mongodb.InitConnection()
	default:
		log.Fatal("Wrong value for usingDatabase parameter, check config")
	}

	handler, _ := handlers.NewClientHandler(repoClient)
	userHandler, _ := handlers.NewUserHandler(repoUsers)
	router := gin.Default()
	router.GET("/clients", middleware.AuthForOperation(redisToken, "clients", "read"), handler.GetClients)
	router.GET("/clients/:id", middleware.AuthForOperation(redisToken, "clients", "read"), handler.GetClientById)
	router.POST("/clients", middleware.AuthForOperation(redisToken, "clients", "create"), handler.CreateClient)
	router.PATCH("/clients", middleware.AuthForOperation(redisToken, "clients", "update"), handler.UpdateClient)
	router.DELETE("/clients/:id", middleware.AuthForOperation(redisToken, "clients", "delete"), handler.DeleteClient)

	router.GET("/users", userHandler.List)
	router.GET("/users/:id", userHandler.ById)
	router.POST("/users", userHandler.CreateUser)
	router.PATCH("/users", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)

	router.POST("/login", middleware.Login(userHandler, redisToken))
	router.POST("/logout", middleware.Logout(redisToken))

	newGraph, err := graph.NewGraph(repoClient)
	if err != nil {
		return
	}

	router.POST("/graph", middleware.Auth(redisToken), newGraph.GraphqlHandler)

	router.Run("127.0.0.1:8080")
}

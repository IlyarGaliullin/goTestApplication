package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"testApplication/interfaces"
	"testApplication/models"
	"testApplication/repositories/postgres"
	"testApplication/utils"
)

var pg interfaces.ClientRepo

func getClients(c *gin.Context) {

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))

	clients := pg.GetClients(c, offset, limit)

	c.IndentedJSON(http.StatusOK, clients)
}

func getClientById(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	client, err := pg.GetClientById(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
	}
	c.IndentedJSON(http.StatusOK, client)
}

func createClient(c *gin.Context) {

	var client models.Client

	err := c.BindJSON(&client)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	lastId, err := pg.CreateClient(c, client)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success", "id": lastId})
}

func updateClient(c *gin.Context) {

	var client models.Client

	err := c.BindJSON(&client)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	err = pg.UpdateClient(c, client)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success"})
}

func deleteClient(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	err := pg.DeleteClient(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success"})
}

func main() {

	utils.LoadConf()

	pg = postgres.InitConnection()
	router := gin.Default()
	router.GET("/clients", getClients)
	router.GET("/clients/:id", getClientById)
	router.POST("/clients", createClient)
	router.PATCH("/clients", updateClient)
	router.DELETE("/clients/:id", deleteClient)

	router.Run("127.0.0.1:8080")

}

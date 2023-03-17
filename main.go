package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"testApplication/databaseDriver"
	"testApplication/models"
	"testApplication/repositories"
)

var db *sql.DB
var pg *repositories.Postgres

func getClients(c *gin.Context) {

	clients := pg.GetClients()

	c.IndentedJSON(http.StatusOK, clients)
}

func getClientById(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	client, err := pg.GetClientById(id)
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

	lastId, err := pg.CreateClient(client)
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error"})
		return
	}

	rowCount, err := pg.UpdateClient(client)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success", "rowCount": rowCount})
}

func deleteClient(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	rowCount, err := pg.DeleteClient(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success", "rowCount": rowCount})
}

func main() {

	db = databaseDriver.NewConnection()
	defer db.Close()

	pg = repositories.NewPostgres(db)
	router := gin.Default()
	router.GET("/clients", getClients)
	router.GET("/clients/:id", getClientById)
	router.POST("/clients", createClient)
	router.PATCH("/clients", updateClient)
	router.DELETE("/clients/:id", deleteClient)

	router.Run("127.0.0.1:8080")

}

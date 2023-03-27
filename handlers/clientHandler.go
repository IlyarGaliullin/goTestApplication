package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"testApplication/interfaces"
	"testApplication/models"
)

type clientHandler struct {
	repo interfaces.ClientRepo
}

func NewClientHandler(repo interfaces.ClientRepo) (*clientHandler, error) {

	clientHandler := clientHandler{
		repo: repo,
	}

	return &clientHandler, nil
}

func (handler *clientHandler) GetClients(c *gin.Context) {

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))

	clients, err := handler.repo.GetClients(c, offset, limit)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	c.IndentedJSON(http.StatusOK, clients)
}

func (handler *clientHandler) GetClientById(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	client, err := handler.repo.GetClientById(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.IndentedJSON(http.StatusOK, client)
}

func (handler *clientHandler) CreateClient(c *gin.Context) {

	var client models.Client

	err := c.BindJSON(&client)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	insertedClient, err := handler.repo.CreateClient(c, client)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success", "client": insertedClient})
}

func (handler *clientHandler) UpdateClient(c *gin.Context) {

	var client models.Client

	err := c.BindJSON(&client)
	if err != nil {

		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	err = handler.repo.UpdateClient(c, client)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success"})
}

func (handler *clientHandler) DeleteClient(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	err := handler.repo.DeleteClient(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success"})
}

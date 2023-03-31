package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"testApplication/interfaces"
	"testApplication/models"
)

type UserHandler struct {
	Repo interfaces.UserRepo
}

func NewUserHandler(repo interfaces.UserRepo) (*UserHandler, error) {

	userHandler := UserHandler{
		Repo: repo,
	}

	return &userHandler, nil
}

func (handler *UserHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))

	user, err := handler.Repo.List(c, offset, limit)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (handler *UserHandler) ById(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	user, err := handler.Repo.ById(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func (handler *UserHandler) CreateUser(c *gin.Context) {

	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	insertedUser, err := handler.Repo.CreateUser(c, user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success", "user": insertedUser})
}

func (handler *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	user, err = handler.Repo.UpdateUser(c, user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success", "user": user})
}

func (handler *UserHandler) DeleteUser(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	user, err := handler.Repo.DeleteUser(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success", "user": user})
}

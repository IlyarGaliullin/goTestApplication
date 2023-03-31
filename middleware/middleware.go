package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"testApplication/handlers"
	"testApplication/interfaces"
	"testApplication/models"
	"testApplication/redis"
	"testApplication/repositories/postgres"
)

func Auth(redisToken *redis.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("token")
		if token == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "empty token"})
			c.Abort()
			return
		}

		_, err := redisToken.CheckToken(c, token)
		if err != nil {
			fmt.Println(err)
			if err == redis.ErrUnauthorized {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{})
				c.Abort()
				return
			}
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Next()
		return
	}
}

func AuthForOperation(redisToken *redis.Connection, table string, operation string) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("token")
		if token == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "empty token"})
			c.Abort()
			return
		}

		userId, err := redisToken.CheckToken(c, token)
		if err != nil {
			fmt.Println(err)
			if err == redis.ErrUnauthorized {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{})
				c.Abort()
				return
			}
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		grants, err := postgres.InitConnection().GetAllUserGrants(c, userId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		for _, v := range grants[table] {
			if v == operation {
				c.Next()
				return
			}
		}
		c.IndentedJSON(http.StatusUnauthorized, gin.H{})
		c.Abort()
		return
	}
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func Login(handler *handlers.UserHandler, redisToken *redis.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.User
		err := c.Bind(&user)

		userFromDb, err := handler.Repo.ByEmail(c, user.Email)

		if err != nil {
			if err == interfaces.ErrNoRows {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Wrong credentials"})
				return
			}
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
		if user.Password != userFromDb.Password {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Wrong credentials"})
			return
		}

		token := generateSecureToken(10)

		err = redisToken.AddToken(c, userFromDb, token)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"token": token})
	}
}

func Logout(redisToken *redis.Connection) gin.HandlerFunc {

	return func(c *gin.Context) {

		token := c.GetHeader("token")
		if token == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty token"})
		}

		err := redisToken.RemoveToken(c, token)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"result": "OK"})
	}
}

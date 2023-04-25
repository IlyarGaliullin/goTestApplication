package middleware

import (
	"crypto"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"testApplication/handlers"
	"testApplication/interfaces"
	"testApplication/models"
	"testApplication/redis"
	"time"
)

func getToken(c *gin.Context) (token string) {
	token = c.GetHeader("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		token = strings.Split(token, " ")[1]
		return
	}
	return
}

func Auth(redisConn *redis.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := getToken(c)
		if token == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "empty token"})
			c.Abort()
			return
		}

		_, err := redisConn.CheckToken(c, token)
		if err != nil {
			log.Println(err)
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

func AuthForOperation(redisConn *redis.Connection, userRepo interfaces.UserRepo, table string, operation string) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := getToken(c)
		log.Printf("authentication attempt from %s, token: %s", c.ClientIP(), token)
		if token == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "empty token"})
			c.Abort()
			return
		}
		userId, err := redisConn.CheckToken(c, token)
		if err != nil {
			log.Println(err)
			if err == redis.ErrUnauthorized {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{})
				log.Printf("authentication failed from %s, token: %s", c.ClientIP(), token)
				c.Abort()
				return
			}
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
			log.Printf("authentication failed from %s, token: %s", c.ClientIP(), token)
			c.Abort()
			return
		}

		grant, err := userRepo.CheckUserGrant(c, userId, table, operation)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
			log.Printf("authentication failed from %s, token: %s, user id: %d", c.ClientIP(), token, userId)
			c.Abort()
			return
		}
		if grant {
			log.Printf("authentication successfull from %s, token: %s, user id: %d", c.ClientIP(), token, userId)
			c.Next()
			return
		}
		c.IndentedJSON(http.StatusUnauthorized, gin.H{})
		c.Abort()
		log.Printf("authentication failed from %s, token: %s, user id: %d", c.ClientIP(), token, userId)
		return
	}
}

func generateSecureToken(email string) string {
	sha := crypto.SHA256.New()
	secureString := fmt.Sprint(email, time.Now().Unix())
	sha.Write([]byte(secureString))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func Login(userHandler *handlers.UserHandler, redisConn *redis.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.User
		err := c.Bind(&user)
		if err != nil {
			log.Printf("authorization failed from %s, email: %s", c.ClientIP(), user.Email)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong parameters"})
			return
		}

		userFromDb, err := userHandler.Repo.ByEmail(c, user.Email)
		log.Printf("authorization attempt from %s, email: %s", c.ClientIP(), user.Email)

		if err != nil {
			if err == interfaces.ErrNoRows {
				log.Printf("authorization failed from %s, email: %s", c.ClientIP(), user.Email)
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Wrong credentials"})
				return
			}
			log.Printf("authorization failed from %s, email: %s, internal server error %s", c.ClientIP(), user.Email, err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		if compare := bcrypt.CompareHashAndPassword([]byte(userFromDb.Password), []byte(user.Password)); compare != nil {
			log.Printf("authorization failed from %s, email: %s", c.ClientIP(), user.Email)
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Wrong credentials"})
			return
		}

		token := generateSecureToken(user.Email)

		err = redisConn.AddToken(c, userFromDb, token)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
			log.Printf("authorization failed from %s, email: %s, internal server error %s", c.ClientIP(), user.Email, err)
			return
		}
		log.Printf("authorization success from %s, email: %s, token: %s", c.ClientIP(), user.Email, token)
		c.IndentedJSON(http.StatusOK, gin.H{"token": token})
	}
}

func Logout(redisConn *redis.Connection) gin.HandlerFunc {

	return func(c *gin.Context) {

		token := getToken(c)
		if token == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "empty token"})
		}

		err := redisConn.RemoveToken(c, token)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"result": "OK"})
	}
}

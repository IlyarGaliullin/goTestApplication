package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testApplication/handlers"
	"testApplication/redis"
	"testApplication/repositories/postgres"
	"testApplication/utils"
	"testing"
)

func TestLogin(t *testing.T) {

	utils.LoadConf()

	pg := postgres.InitConnectionNoMigration()
	userHandler, err := handlers.NewUserHandler(pg)
	redisConn, err := redis.NewConn()
	if err != nil {
		panic(err)
	}

	router := gin.New()

	router.POST("/login", Login(userHandler, redisConn))

	jsonBody := []byte(`{"email":"admin@admin.adm","password":"admin"}`)

	fmt.Println(string(jsonBody))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
}

func TestFullLoginCycle(t *testing.T) {

	utils.LoadConf()

	pg := postgres.InitConnectionNoMigration()
	userHandler, err := handlers.NewUserHandler(pg)
	redisConn, err := redis.NewConn()
	if err != nil {
		panic(err)
	}

	router := gin.New()

	router.POST("/login", Login(userHandler, redisConn))

	jsonBody := []byte(`{"email":"admin@admin.adm","password":"admin"}`)

	fmt.Println(string(jsonBody))
	recorderLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
	reqLogin.Header.Add("Content-Type", "application/json")
	if err != nil {
		return
	}
	router.ServeHTTP(recorderLogin, reqLogin)

	assert.Equal(t, recorderLogin.Code, http.StatusOK)
	body, err := io.ReadAll(recorderLogin.Body)
	if err != nil {
		t.Fatal(err)
	}
	var loginResponse struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		t.Fatal(err)
	}

	router.GET("/auth", Auth(redisConn))
	recorderAuth := httptest.NewRecorder()
	reqAuth, _ := http.NewRequest("GET", "/auth", nil)
	reqAuth.Header.Add("Content-Type", "application/json")
	reqAuth.Header.Add("Authorization", "Bearer "+loginResponse.Token)

	router.ServeHTTP(recorderAuth, reqAuth)
	assert.Equal(t, recorderAuth.Code, http.StatusOK)

	router.POST("/logout", Logout(redisConn))
	recorderLogout := httptest.NewRecorder()
	reqLogout, _ := http.NewRequest("POST", "/logout", nil)
	reqLogout.Header.Add("Content-Type", "application/json")
	reqLogout.Header.Add("Authorization", loginResponse.Token)

	router.ServeHTTP(recorderLogout, reqLogout)
	assert.Equal(t, recorderLogout.Code, http.StatusOK)
	fmt.Println(recorderLogout.Body.String())
}

func TestLogin401(t *testing.T) {

	utils.LoadConf()

	pg := postgres.InitConnectionNoMigration()
	userHandler, err := handlers.NewUserHandler(pg)
	redisConn, err := redis.NewConn()
	if err != nil {
		panic(err)
	}

	router := gin.New()

	router.POST("/login", Login(userHandler, redisConn))

	jsonBody := []byte(`{"email":"nonexisting","password":"wrong"}`)

	fmt.Println(string(jsonBody))
	recorderLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
	reqLogin.Header.Add("Content-Type", "application/json")
	if err != nil {
		return
	}
	router.ServeHTTP(recorderLogin, reqLogin)

	want, _ := json.MarshalIndent(
		struct {
			Message string `json:"message"`
		}{
			Message: "Wrong credentials",
		}, "", "    ")
	assert.Equal(t, recorderLogin.Code, http.StatusUnauthorized)
	assert.Equal(t, recorderLogin.Body.String(), string(want))
}

func TestAuthForOperation(t *testing.T) {

	utils.LoadConf()

	pg := postgres.InitConnectionNoMigration()
	userHandler, err := handlers.NewUserHandler(pg)
	redisConn, err := redis.NewConn()
	if err != nil {
		panic(err)
	}

	router := gin.New()

	router.POST("/login", Login(userHandler, redisConn))

	jsonBody := []byte(`{"email":"admin@admin.adm","password":"admin"}`)

	fmt.Println(string(jsonBody))
	recorderLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
	reqLogin.Header.Add("Content-Type", "application/json")
	if err != nil {
		return
	}
	router.ServeHTTP(recorderLogin, reqLogin)

	assert.Equal(t, recorderLogin.Code, http.StatusOK)
	body, err := io.ReadAll(recorderLogin.Body)
	if err != nil {
		t.Fatal(err)
	}
	var loginResponse struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		t.Fatal(err)
	}

	router.GET("/authOperation", AuthForOperation(redisConn, pg, "clients", "read"))
	recorderAuth := httptest.NewRecorder()
	reqAuth, _ := http.NewRequest("GET", "/authOperation", nil)
	reqAuth.Header.Add("Content-Type", "application/json")
	reqAuth.Header.Add("Authorization", "Bearer "+loginResponse.Token)

	router.ServeHTTP(recorderAuth, reqAuth)
	assert.Equal(t, recorderAuth.Code, http.StatusOK)
}

func TestAuthForOperation401(t *testing.T) {

	utils.LoadConf()

	pg := postgres.InitConnectionNoMigration()
	userHandler, err := handlers.NewUserHandler(pg)
	redisConn, err := redis.NewConn()
	if err != nil {
		panic(err)
	}

	router := gin.New()

	router.POST("/login", Login(userHandler, redisConn))

	jsonBody := []byte(`{"email":"user@user.usr","password":"user"}`)

	fmt.Println(string(jsonBody))
	recorderLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
	reqLogin.Header.Add("Content-Type", "application/json")
	if err != nil {
		return
	}
	router.ServeHTTP(recorderLogin, reqLogin)

	assert.Equal(t, recorderLogin.Code, http.StatusOK)
	body, err := io.ReadAll(recorderLogin.Body)
	if err != nil {
		t.Fatal(err)
	}
	var loginResponse struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		t.Fatal(err)
	}

	router.GET("/authOperation", AuthForOperation(redisConn, pg, "clients", "delete"))
	recorderAuth := httptest.NewRecorder()
	reqAuth, _ := http.NewRequest("GET", "/authOperation", nil)
	reqAuth.Header.Add("Content-Type", "application/json")
	reqAuth.Header.Add("Authorization", "Bearer "+loginResponse.Token)

	router.ServeHTTP(recorderAuth, reqAuth)
	assert.Equal(t, recorderAuth.Code, http.StatusUnauthorized)
}

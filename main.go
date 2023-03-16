package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"testApplication/databaseDriver"
)

var db *sql.DB

//TODO вынести client и crud в отдельный package, создать отдельный пакейдж-хэндлер

type Client struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func getClients(c *gin.Context) {

	var clients []Client

	clientsStmt, err := db.Prepare("select * from clients")
	if err != nil {
		log.Fatal(err)
	}

	defer clientsStmt.Close()
	rows, err := clientsStmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {

		var (
			id   int
			name string
		)

		err := rows.Scan(&id, &name)
		clients = append(clients, Client{
			Id:   id,
			Name: name,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, clients)
}

func getClientById(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	var name string

	clientByIdStmt, err := db.Prepare("select name from clients where id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer clientByIdStmt.Close()

	err = clientByIdStmt.QueryRow(id).Scan(&name)
	if err != nil {

		if err == sql.ErrNoRows {
			fmt.Printf("No client found by id %d\n", id)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No client found by id %d", id)})
			return
		} else {
			log.Fatal(err)
		}
	}
	fmt.Println(id, name)
	c.IndentedJSON(http.StatusOK, Client{Id: id, Name: name})
}

func createClient(c *gin.Context) {

	var newClient Client

	err := c.BindJSON(&newClient)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error"})
		return
	}

	insertClientStmt, err := db.Prepare("insert into clients(name) values($1) returning id")
	if err != nil {
		log.Fatal(err)
	}
	defer insertClientStmt.Close()

	lastId := 0
	err = insertClientStmt.QueryRow(newClient.Name).Scan(&lastId)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("1 row inserted, id = %d\n", lastId)

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success", "id": lastId})
}

func updateClient(c *gin.Context) {

	var client Client

	err := c.BindJSON(&client)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Error"})
		return
	}

	updateClientStmt, err := db.Prepare("update clients set name = $1 where id = $2")
	if err != nil {
		log.Fatal(err)
	}
	defer updateClientStmt.Close()

	res, err := updateClientStmt.Exec(client.Name, client.Id)
	if err != nil {
		log.Fatal(err)
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success", "rowCount": rowCount})
}

func deleteClient(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	insertClientStmt, err := db.Prepare("delete from clients where id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer insertClientStmt.Close()

	res, err := insertClientStmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	rowCount, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Success", "rowCount": rowCount})
}

func main() {

	db = databaseDriver.NewConnection()
	defer db.Close()

	router := gin.Default()
	router.GET("/clients", getClients)
	router.GET("/clients/:id", getClientById)
	router.POST("/clients", createClient)
	router.PATCH("/clients", updateClient)
	router.DELETE("/clients/:id", deleteClient)

	router.Run("127.0.0.1:8080")

}

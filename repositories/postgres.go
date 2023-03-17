package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testApplication/models"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (pg *Postgres) GetClients() []models.Client {
	var clients []models.Client

	clientsStmt, err := pg.db.Prepare("select * from clients")
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
		clients = append(clients, models.Client{
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

	return clients
}

func (pg *Postgres) GetClientById(id int) (models.Client, error) {

	var name string

	clientByIdStmt, err := pg.db.Prepare("select name from clients where id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer clientByIdStmt.Close()

	err = clientByIdStmt.QueryRow(id).Scan(&name)
	if err != nil {

		if err == sql.ErrNoRows {
			fmt.Printf("No client found by id %d\n", id)
			return models.Client{}, errors.New(fmt.Sprintf("No client found by id %d\n", id))
		} else {
			log.Fatal(err)
		}
	}
	return models.Client{Id: id, Name: name}, nil
}

func (pg *Postgres) CreateClient(newClient models.Client) (int, error) {

	insertClientStmt, err := pg.db.Prepare("insert into clients(name) values($1) returning id")
	if err != nil {
		return 0, err
	}
	defer insertClientStmt.Close()

	lastId := 0
	err = insertClientStmt.QueryRow(newClient.Name).Scan(&lastId)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("1 row inserted, id = %d\n", lastId)
	return lastId, nil
}

func (pg *Postgres) UpdateClient(client models.Client) (int, error) {

	updateClientStmt, err := pg.db.Prepare("update clients set name = $1 where id = $2")
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

	fmt.Printf("%d row affected\n", rowCount)

	return int(rowCount), nil
}

func (pg *Postgres) DeleteClient(id int) (int, error) {

	insertClientStmt, err := pg.db.Prepare("delete from clients where id = $1")
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
	fmt.Printf("%d row deleted\n", rowCount)
	return int(rowCount), nil
}

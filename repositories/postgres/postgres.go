package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testApplication/models"
	"testApplication/utils"
)

type postgres struct {
	db *sql.DB
}

func InitConnection() *postgres {

	host := utils.Conf.Get("postgres.host")
	port := utils.Conf.GetInt("postgres.port")
	user := utils.Conf.Get("postgres.user")
	password := utils.Conf.Get("postgres.password")
	database := utils.Conf.Get("postgres.database")

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)

	var err error
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &postgres{db: db}
}

func (pg *postgres) GetClients(ctx context.Context, offset int, limit int) []models.Client {
	var clients []models.Client

	clientsStmt, err := pg.db.Prepare("select * from clients order by id limit $2 offset $1")
	if err != nil {
		log.Fatal(err)
	}
	defer clientsStmt.Close()

	var rows *sql.Rows
	if limit == 0 {
		rows, err = clientsStmt.Query(offset, nil)
	} else {
		rows, err = clientsStmt.Query(offset, limit)
	}
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

func (pg *postgres) GetClientById(ctx context.Context, id int) (models.Client, error) {

	var name string

	clientByIdStmt, err := pg.db.Prepare("select name from clients where id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer clientByIdStmt.Close()

	err = clientByIdStmt.QueryRow(id).Scan(&name)
	if err != nil {

		if err == sql.ErrNoRows {
			return models.Client{}, errors.New(fmt.Sprintf("No client found by id %d\n", id))
		} else {
			log.Fatal(err)
		}
	}
	return models.Client{Id: id, Name: name}, nil
}

func (pg *postgres) CreateClient(ctx context.Context, newClient models.Client) (models.Client, error) {

	insertClientStmt, err := pg.db.Prepare("insert into clients(name) values($1) returning *")
	if err != nil {
		return models.Client{}, err
	}
	defer insertClientStmt.Close()

	lastId := 0
	name := ""
	err = insertClientStmt.QueryRow(newClient.Name).Scan(&lastId, &name)
	if err != nil {
		log.Fatal(err)
	}

	return models.Client{Id: lastId, Name: name}, nil
}

func (pg *postgres) UpdateClient(ctx context.Context, client models.Client) error {

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
	if rowCount == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (pg *postgres) DeleteClient(ctx context.Context, id int) error {

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
	if rowCount == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

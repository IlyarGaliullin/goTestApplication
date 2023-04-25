package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"testApplication/interfaces"
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

	driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{})
	migr, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = migr.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &postgres{db: db}
}

func InitConnectionNoMigration() *postgres {

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

func (pg *postgres) GetClients(ctx context.Context, offset int, limit int) ([]models.Client, error) {
	var clients []models.Client

	clientsStmt, err := pg.db.Prepare("SELECT * FROM clients ORDER BY id LIMIT $2 OFFSET $1")
	if err != nil {
		log.Println(err)
		return clients, err
	}
	defer clientsStmt.Close()

	var rows *sql.Rows
	if limit == 0 {
		rows, err = clientsStmt.Query(offset, nil)
	} else {
		rows, err = clientsStmt.Query(offset, limit)
	}
	if err != nil {
		log.Println(err)
		return clients, err
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
			log.Println(err)
			return clients, err
		}
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return clients, err
	}

	return clients, nil
}

func (pg *postgres) GetClientById(ctx context.Context, id int) (models.Client, error) {

	var name string

	clientByIdStmt, err := pg.db.Prepare("SELECT name FROM clients WHERE id = $1")
	if err != nil {
		log.Println(err)
		return models.Client{}, err
	}
	defer clientByIdStmt.Close()

	err = clientByIdStmt.QueryRow(id).Scan(&name)
	if err != nil {

		if err == sql.ErrNoRows {
			return models.Client{}, errors.New(fmt.Sprintf("No client found by id %d\n", id))
		} else {
			log.Println(err)
			return models.Client{}, err
		}
	}
	return models.Client{Id: id, Name: name}, nil
}

func (pg *postgres) CreateClient(ctx context.Context, newClient models.Client) (models.Client, error) {

	insertClientStmt, err := pg.db.Prepare("INSERT INTO clients(name) VALUES($1) returning *")
	if err != nil {
		return models.Client{}, err
	}
	defer insertClientStmt.Close()

	lastId := 0
	name := ""
	err = insertClientStmt.QueryRow(newClient.Name).Scan(&lastId, &name)
	if err != nil {
		log.Println(err)
		return models.Client{}, err
	}

	return models.Client{Id: lastId, Name: name}, nil
}

func (pg *postgres) UpdateClient(ctx context.Context, client models.Client) error {

	updateClientStmt, err := pg.db.Prepare("UPDATE clients SET name = $1 WHERE id = $2")
	if err != nil {
		log.Println(err)
		return err
	}
	defer updateClientStmt.Close()

	res, err := updateClientStmt.Exec(client.Name, client.Id)
	if err != nil {
		log.Println(err)
		return err
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}
	if rowCount == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (pg *postgres) DeleteClient(ctx context.Context, id int) error {

	deleteClientStmt, err := pg.db.Prepare("DELETE FROM clients WHERE id = $1")
	if err != nil {
		log.Println(err)
		return err
	}
	defer deleteClientStmt.Close()

	res, err := deleteClientStmt.Exec(id)
	if err != nil {
		log.Println(err)
		return err
	}
	rowCount, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}
	if rowCount == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (pg *postgres) List(ctx context.Context, offset int, limit int) ([]models.User, error) {
	var users []models.User

	usersStmt, err := pg.db.Prepare("SELECT * FROM users ORDER BY id LIMIT $2 OFFSET $1")
	if err != nil {
		log.Println(err)
		return users, err
	}
	defer usersStmt.Close()

	var rows *sql.Rows
	if limit == 0 {
		rows, err = usersStmt.Query(offset, nil)
	} else {
		rows, err = usersStmt.Query(offset, limit)
	}
	if err != nil {
		log.Println(err)
		return users, err
	}

	for rows.Next() {

		var (
			id   int
			name string
		)

		err := rows.Scan(&id, &name)
		users = append(users, models.User{
			Id:   id,
			Name: name,
		})
		if err != nil {
			log.Println(err)
			return users, err
		}
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return users, err
	}

	return users, nil
}

func (pg *postgres) ById(ctx context.Context, id int) (models.User, error) {

	var name string

	userByIdStmt, err := pg.db.Prepare("SELECT name FROM user WHERE id = $1")
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}
	defer userByIdStmt.Close()

	err = userByIdStmt.QueryRow(id).Scan(&name)
	if err != nil {

		if err == sql.ErrNoRows {
			return models.User{}, errors.New(fmt.Sprintf("No user found by id %d", id))
		} else {
			log.Println(err)
			return models.User{}, err
		}
	}
	return models.User{Id: id, Name: name}, nil
}

func (pg *postgres) ByEmail(ctx context.Context, email string) (models.User, error) {

	var id int
	var password string

	userByEmailStmt, err := pg.db.Prepare("SELECT id, password FROM users WHERE email = $1")
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}
	defer userByEmailStmt.Close()

	err = userByEmailStmt.QueryRow(email).Scan(&id, &password)

	if err != nil {

		if err == sql.ErrNoRows {
			return models.User{}, interfaces.ErrNoRows
		} else {
			log.Println(err)
			return models.User{}, err
		}
	}
	return models.User{Id: id, Password: password}, nil
}

func (pg *postgres) CreateUser(ctx context.Context, newUser models.User) (models.User, error) {

	insertUserStmt, err := pg.db.Prepare("INSERT INTO users(name, email, password) VALUES($1, $2, $3) returning id, name, email")
	if err != nil {
		return models.User{}, err
	}
	defer insertUserStmt.Close()

	lastId := 0
	name := ""
	email := ""
	err = insertUserStmt.QueryRow(newUser.Name, newUser.Email, newUser.Password).Scan(&lastId, &name, &email)
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}

	return models.User{Id: lastId, Name: name, Email: email}, nil
}

func (pg *postgres) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	updateClientStmt, err := pg.db.Prepare("UPDATE users SET name = $1  WHERE id = $2")
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}
	defer updateClientStmt.Close()

	res, err := updateClientStmt.Exec(user.Name, user.Id)
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}
	if rowCount == 0 {
		return models.User{}, errors.New("no rows affected")
	}

	userUpdated, _ := pg.ById(ctx, user.Id)
	return userUpdated, nil
}

func (pg *postgres) DeleteUser(ctx context.Context, id int) (models.User, error) {

	deleteUserStmt, err := pg.db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}
	defer deleteUserStmt.Close()

	user, err := pg.ById(ctx, id)

	res, err := deleteUserStmt.Exec(id)
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}
	rowCount, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}
	if rowCount == 0 {
		return models.User{}, errors.New("no rows affected")
	}

	return user, nil
}

func (pg *postgres) UpdateRoles(ctx context.Context, user models.User, roles []models.Role) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (pg *postgres) GetAllUserGrants(ctx context.Context, userId int) (grants []models.Grant, err error) {

	grantsStmt, err := pg.db.Prepare(
		"SELECT g.ontable, g.read, g.\"create\", g.update, g.delete FROM userroles ur " +
			" JOIN roles r ON ur.roleid = r.id" +
			" JOIN grants g on r.id = g.roleid" +
			" WHERE ur.userid = $1",
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer grantsStmt.Close()

	var rows *sql.Rows
	rows, err = grantsStmt.Query(userId)
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		var grant models.Grant
		err = rows.Scan(&grant.Table, &grant.Read, &grant.Create, &grant.Update, &grant.Delete)
		var operations []string

		if grant.Read {
			operations = append(operations, "read")
		}
		if grant.Create {
			operations = append(operations, "create")
		}
		if grant.Update {
			operations = append(operations, "update")
		}
		if grant.Delete {
			operations = append(operations, "delete")
		}

		grants = append(grants, grant)
		if err != nil {
			log.Println(err)
			return
		}
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (pg *postgres) CheckUserGrant(ctx context.Context, userId int, table string, operation string) (found bool, err error) {

	sqlString := "SELECT true found FROM userroles ur " +
		" JOIN roles r ON ur.roleid = r.id" +
		" JOIN grants g on r.id = g.roleid" +
		" WHERE ur.userid = $1 AND g.ontable = $2"

	sqlString += " AND \"" + operation + "\" = true"

	grantStmt, err := pg.db.Prepare(sqlString)

	if err != nil {
		log.Println(err)
		return false, nil
	}
	defer grantStmt.Close()

	err = grantStmt.QueryRow(userId, table).Scan(&found)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	return found, nil
}

func (pg *postgres) GetFieldsPermissions(ctx context.Context, userId int, tableName string) (map[string]bool, error) {
	//TODO implement me
	panic("implement me")
}

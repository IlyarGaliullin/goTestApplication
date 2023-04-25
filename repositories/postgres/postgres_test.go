package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testApplication/utils"
	"testing"
)

func initPostgres() *postgres {
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

	return &postgres{db: db}
}

func TestPostgres_GetClients(t *testing.T) {

	utils.LoadConf()
	pg := initPostgres()
	t.Run("test", func(t *testing.T) {
		_, err := pg.GetClients(context.TODO(), 0, 2)
		if err != nil {
			return
		}
	})
}

func BenchmarkPostgres_GetClients(b *testing.B) {

	utils.LoadConf()
	pg := initPostgres()
	for i := 0; i < b.N; i++ {
		_, err := pg.GetClients(context.TODO(), 0, 0)
		if err != nil {
			return
		}
	}
}

func TestPostgres_CheckUserGrant(t *testing.T) {

	utils.LoadConf()
	pg := initPostgres()

	var tests = []struct {
		userId    int
		table     string
		operation string
		want      bool
	}{
		{1, "clients", "read", true},
		{3, "clients", "delete", false},
	}
	for _, tt := range tests {
		found, err := pg.CheckUserGrant(context.TODO(), tt.userId, tt.table, tt.operation)
		if err != nil {
			t.Error(err)
		}
		if found != tt.want {
			t.Errorf("got %t, want %t", found, tt.want)
		}
	}
}

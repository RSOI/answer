package database

import (
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/RSOI/answer/utils"
	"github.com/jackc/pgx"
)

var (
	// HOST postgres host
	HOST = "localhost"
	// PORT postgres port
	PORT uint16 = 5432
)

// Connect to postgrss
func Connect() *pgx.ConnPool {
	utils.LOG(fmt.Sprintf("Connecting postgress: %s:%d", HOST, PORT))
	runtime.GOMAXPROCS(runtime.NumCPU())
	connection := pgx.ConnConfig{
		Host:     HOST,
		User:     "dzaytsev",
		Password: "126126",
		Database: "rsoi",
		Port:     PORT,
	}

	var err error
	db, err := pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: connection, MaxConnections: 50})
	if err != nil {
		panic(err)
	}

	err = createShema(db)
	if err != nil {
		panic(err)
	}

	return db
}

func createShema(db *pgx.ConnPool) error {
	utils.LOG("Creating scheme...")

	sql, err := ioutil.ReadFile("database/scheme.sql")
	if err != nil {
		utils.LOG(fmt.Sprintf("Error while creating scheme: %s", err.Error()))
		return err
	}
	shema := string(sql)

	_, err = db.Exec(shema)
	return err
}

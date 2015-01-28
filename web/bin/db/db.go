package db

import (
	"../config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	db     *sql.DB
	logger *log.Logger
)

func init() {

	conf := config.GetConfig()

	file, err := os.OpenFile(`db_errors.log`, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	logger = log.New(file, ``, log.Ldate|log.Ltime|log.Lshortfile)

	db, err = sql.Open("postgres", conf.Storage.Postgresql)
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {

		fmt.Println(err.Error())
		os.Exit(1)
	}

	db.SetMaxIdleConns(8)
	db.SetMaxOpenConns(64)
}

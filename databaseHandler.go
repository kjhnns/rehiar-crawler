package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var dbConnection *sql.DB

func DbConn() *sql.DB {
	return dbConnection
}

func DbClose() {
	dbConnection.Close()
}

func InitDatabase() {
	var err error
	dbConnection, err = sql.Open("postgres", Configuration.DatabaseUrl)
	pingErr := dbConnection.Ping()
	if err != nil || pingErr != nil {
		Configuration.Logger.Error.Println("Database Connection Error: ", err, pingErr)
		panic("Database Configuration")
	}
	Configuration.Logger.Info.Println("Database Connection established")

	_, err = dbConnection.Exec(amazonDataTable)
	if err != nil {
		Configuration.Logger.Error.Println("Creating Table failed: ", err)
	}
}

package main

import (
	"database/sql"
	"fmt"
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
		fmt.Println("\t- Database Connection Error: ", err, pingErr)
		panic("Database Configuration")
	}
	fmt.Println("\t- Database Connection established")

	_, err = dbConnection.Exec(amazonDataTable)
	_, err = dbConnection.Exec(googleTrendsTable)

	if err != nil {
		fmt.Println("\t- Creating Table failed: ", err)
	}
}

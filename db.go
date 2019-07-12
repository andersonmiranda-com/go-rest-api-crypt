package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DBDriver   = "mysql"
	DBName     = "valentium"
	DBUser     = "root"
	DBPassword = "root"
	DBURL      = DBUser + ":" + DBPassword + "@tcp(127.0.0.1:3306)/" + DBName

	MongoDBUri = "mongodb://localhost:27017"
)

/*
const (
	DBDriver   = "mysql"
	DBName     = "atalan_api_db1"
	DBUser     = "atalandb1"
	DBPassword = "KD3uV8>GaHbMbC"
	DBURL      = DBUser + ":" + DBPassword + "@tcp(aa1ozx7xopaip9a.cosyme9gqndq.us-east-1.rds.amazonaws.com:3306)/" + DBName

	MongoDBUri = "mongodb+srv://atalan:mXyD6H-cut8Gp3U@cluster0-vuqrp.mongodb.net"
)
*/

func dbConn() (db *sql.DB) {
	db, err := sql.Open(DBDriver, DBURL)
	if err != nil {
		panic(err.Error())
	}
	return db
}

package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	// import _ "github.com/jinzhu/gorm/dialects/postgres"
	// import _ "github.com/jinzhu/gorm/dialects/sqlite"
	// import _ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"os"
)

var db *gorm.DB //database

func init() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")


	//dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
	dbUri := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username,password,dbHost,dbPort, dbName)
	fmt.Println(dbUri)

	//conn, err := gorm.Open("postgres", dbUri)
	conn, err := gorm.Open("mysql", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &Contact{}) //Database migration
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
)


type DB_ENV struct {
	Host, Port, Dialect, Username, Password, DBname string
}

var (
	success, fail int
	dbEnv DB_ENV
)


func init() {
	goPath := os.Getenv("GOPATH")
	godotenv.Load(goPath+"/src/Auth_System_In_Go/.env")
	dbEnv = DB_ENV{
		Host:os.Getenv("DB_HOST"),
		Port:os.Getenv("DB_PORT"),
		Dialect:os.Getenv("DB_DIALECT"),
		Username:os.Getenv("DB_USERNAME"),
		Password:os.Getenv("DB_PASSWORD"),
		DBname:os.Getenv("DB_NAME"),
	}
}


func DBConnect() (*sql.DB, error) {
	db, _ := sql.Open(dbEnv.Dialect, dbEnv.Username+":"+dbEnv.Password+"@tcp("+dbEnv.Host+":"+dbEnv.Port+")/"+dbEnv.DBname+"?parseTime=true")
	return db, nil
}

func main() {
	if !createSeederTable() {
		fmt.Println("Seeding Failed. Internal Server Error.")
		return
	}
	RoleSeeder()
	UserSeeder()

	if success+fail == 0{
		fmt.Println("No New Seeder, Already Up-to-date.")
	} else {
		fmt.Printf("\nTotal New Seeder = %d\nSuccessfully done = %d\nFailed = %d\n", success+fail, success, fail)
	}
}

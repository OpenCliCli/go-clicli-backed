package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sqlx.DB
	err error
)

func init() {

	godotenv.Load(".env")
	host := os.Getenv("MYSQL_HOSTNAME")
	user := os.Getenv("MYSQL_USERNAME")
	dbName := os.Getenv("MYSQL_DATABASE")
	pwd := os.Getenv("MYSQL_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pwd, host, dbName)
	raw, err := sql.Open("mysql", dsn)
	loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	sqldblogger := sqldblogger.OpenDriver(dsn, raw.Driver(), loggerAdapter)
	db = sqlx.NewDb(sqldblogger, "mysql")

	//https://github.com/docker-library/mysql/issues/124
	fmt.Println("DB INIT ==> ", os.Getenv("ENV"), dsn)

	if err != nil {
		panic(err.Error())
	}
}

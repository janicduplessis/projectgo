package infrastructure

import (
	"github.com/janicduplessis/projectgo/ct/interfaces"

	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlHandler struct {
	Conn *sql.DB
}

type MySqlDbConfig struct {
	User     string
	Password string
	Name     string
	Url      string
	Port     string
}

func NewMySqlHandler(config MySqlDbConfig) *MySqlHandler {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password, config.Url, config.Port, config.Name))
	if err != nil {
		panic(err)
	}

	interfaces.ErrNoRows = sql.ErrNoRows
	conn.SetMaxOpenConns(10)

	return &MySqlHandler{
		Conn: conn,
	}
}

func (handler *MySqlHandler) Execute(statement string, params ...interface{}) (interfaces.Result, error) {
	return handler.Conn.Exec(statement, params...)
}

func (handler *MySqlHandler) Query(statement string, params ...interface{}) (interfaces.Rows, error) {
	return handler.Conn.Query(statement, params...)
}

func (handler *MySqlHandler) QueryRow(statement string, params ...interface{}) interfaces.Row {
	return handler.Conn.QueryRow(statement, params...)
}

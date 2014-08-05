package ct

import (
	"database/sql"

	"code.google.com/p/go.crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"
)

type Auth struct {
	db *sql.DB
}

type RegisterInfo struct {
	UserName  string
	Password  string
	FirstName string
	LastName  string
	Email     string
}

type LoginInfo struct {
	UserName string
	Password string
}

func NewAuth(db *sql.DB) *Auth {
	return &Auth{
		db: db,
	}
}

func (a *Auth) Register(info *RegisterInfo) (bool, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	//Validate infos

	stmt, err := a.db.Prepare("CALL pRegister (?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(info.UserName, string(hash), info.FirstName, info.LastName, info.Email)
	if err != nil {
		panic(err)
	}

	return true, nil
}

func (a *Auth) Login(info *LoginInfo) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	//Validate infos

	stmt, err := a.db.Prepare("CALL pLogin (?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(info.UserName, string(hash))
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		return &User{}, nil
	} else {
		return nil, nil
	}
}

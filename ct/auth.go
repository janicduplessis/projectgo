package ct

import (
	"database/sql"
	//"log"

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
		return false, err
	}
	//Validate infos

	_, err = a.db.Exec("CALL pRegister (?, ?, ?, ?, ?)", info.UserName, string(hash), info.FirstName, info.LastName, info.Email)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *Auth) Login(info *LoginInfo) (*User, error) {
	//Validate infos
	rows, err := a.db.Query(`SELECT UserId, PasswordHash, UserName, FirstName, LastName, Email
							 FROM user
							 WHERE UserName = ?`,
		info.UserName)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var (
			id        int
			hash      string
			username  string
			firstName string
			lastName  string
			email     string
		)
		err = rows.Scan(&id, &hash, &username, &firstName, &lastName, &email)
		if err != nil {
			return nil, err
		}
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(info.Password))
		if err != nil {
			return nil, nil
		}

		return &User{
			UserId:    id,
			UserName:  username,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		}, nil

	} else {
		return nil, nil
	}
}

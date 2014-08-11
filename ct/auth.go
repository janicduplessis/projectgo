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

func (a *Auth) Register(info *RegisterInfo) (*User, error) {

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// TODO: Validate infos

	// Check if the username is available
	err = a.db.QueryRow(`SELECT 1
				 	   FROM user
				 	   WHERE UserName = ?`, info.UserName).Scan(new(int))
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, err
		}

		// Username unavailable
		return nil, nil
	}

	// Register the user!
	res, err := a.db.Exec(`INSERT INTO user (UserName, PasswordHash, FirstName, LastName, Email)
						   VALUES (?, ?, ?, ?, ?)`,
		info.UserName, string(hash), info.FirstName, info.LastName, info.Email)

	if err != nil {
		return nil, err
	}

	userId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{
		UserId:    userId,
		UserName:  info.UserName,
		FirstName: info.FirstName,
		LastName:  info.LastName,
		Email:     info.Email,
	}, nil
}

func (a *Auth) Login(info *LoginInfo) (*User, error) {

	// Validate infos

	var (
		id        int64
		hash      string
		username  string
		firstName string
		lastName  string
		email     string

		user *User
	)

	// Returns info for the UserName
	err := a.db.QueryRow(`SELECT UserId, PasswordHash, UserName, FirstName, LastName, Email
						  FROM user
						  WHERE UserName = ?`,
		info.UserName).Scan(&id, &hash, &username, &firstName, &lastName, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Bad username
			return nil, nil
		}

		return nil, err
	}

	// Validate the password
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(info.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			// Bad password
			return nil, nil
		}

		return nil, err
	}

	// All good, create the user object
	user = &User{
		UserId:    id,
		UserName:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	return user, nil
}

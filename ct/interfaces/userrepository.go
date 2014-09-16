package interfaces

import (
	"github.com/janicduplessis/projectgo/ct/domain"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

type DbUserRepo DbRepo

func NewDbUserRepo(dbHandlers map[string]DbHandler) *DbUserRepo {
	dbHandler := dbHandlers["DbUserRepo"]

	return &DbUserRepo{
		dbHandlers: dbHandlers,
		dbHandler:  dbHandler,
	}
}

func (repo *DbUserRepo) Create(info *usecases.RegisterInfo) (*usecases.User, error) {
	// Check if the username is available
	err := repo.dbHandler.QueryRow(`SELECT 1
				 	    		    FROM user
				 	     		    WHERE Username = ?`, info.Username).Scan(new(int))

	if err != ErrNoRows {
		if err == nil {
			return nil, usecases.ErrUserAlreadyExists
		}
		return nil, err
	}

	// Register the user!
	// Insert in user
	res, err := repo.dbHandler.Execute(`INSERT INTO user (Username, PasswordHash)
						   			    VALUES (?, ?)`,
		info.Username, info.Password)

	if err != nil {
		return nil, err
	}

	userId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Insert in client
	res, err = repo.dbHandler.Execute(`INSERT INTO client (ClientId, DisplayName, FirstName, LastName, Email)
						   			   VALUES (?, ?, ?, ?, ?)`, userId, info.Username, info.FirstName, info.LastName, info.Email)

	if err != nil {
		return nil, err
	}

	// Create objects
	client := &domain.Client{
		Id:          userId,
		DisplayName: info.Username,
		FirstName:   info.FirstName,
		LastName:    info.LastName,
		Email:       info.Email,
	}

	user := &usecases.User{
		Id:       userId,
		Username: info.Username,
		Client:   client,
	}

	return user, nil
}

func (repo *DbUserRepo) FindByName(name string) (*usecases.User, error) {
	user, _, err := repo.FindByNameWithHash(name)
	return user, err
}

func (repo *DbUserRepo) FindByNameWithHash(name string) (*usecases.User, string, error) {
	var (
		userId       int64
		username     string
		passwordHash string
		displayName  string
		firstName    string
		lastName     string
		email        string
	)
	err := repo.dbHandler.QueryRow(`SELECT u.UserId, u.Username, u.PasswordHash, c.DisplayName, c.FirstName, c.LastName, c.Email
				 	    		    FROM user u
				 	    		    JOIN client c ON u.UserId = c.ClientId
				 	     		    WHERE u.Username = ?`, name).Scan(&userId, &username, &passwordHash, &displayName, &firstName, &lastName, &email)

	if err != nil {
		if err == ErrNoRows {
			return nil, "", nil
		}
		return nil, "", err
	}

	// Create objects
	client := &domain.Client{
		Id:          userId,
		DisplayName: displayName,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
	}

	user := &usecases.User{
		Id:       userId,
		Username: username,
		Client:   client,
	}

	return user, passwordHash, nil
}

func (repo *DbUserRepo) UpdatePasswordHash(user *usecases.User, hash string) error {
	_, err := repo.dbHandler.Execute(`UPDATE user SET PasswordHash=?
									  WHERE UserId=?`, hash, user.Id)

	return err
}

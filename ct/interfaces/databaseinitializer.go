package interfaces

type DbInitializerRepo DbRepo

func NewDbInitializerRepo(dbHandlers map[string]DbHandler) *DbInitializerRepo {
	dbHandler := dbHandlers["DbInitializerRepo"]

	return &DbInitializerRepo{
		dbHandlers: dbHandlers,
		dbHandler:  dbHandler,
	}
}

func (repo *DbInitializerRepo) Init() {
	// Create the tables if it doesnt exists
	_, err := repo.dbHandler.Execute(`CREATE TABLE IF NOT EXISTS client (
								  ClientId int(11) NOT NULL AUTO_INCREMENT,
								  FirstName varchar(255) NOT NULL,
								  LastName varchar(255) NOT NULL,
								  Email varchar(255) NOT NULL,
								  PRIMARY KEY (ClientId)
							  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;`)

	if err != nil {
		panic(err)
	}

	_, err = repo.dbHandler.Execute(`CREATE TABLE IF NOT EXISTS user (
								  UserId int(11) NOT NULL AUTO_INCREMENT,
								  Username varchar(64) NOT NULL,
								  PasswordHash varchar(64) NOT NULL,
								  ClientId int(11) NOT NULL,
								  PRIMARY KEY (UserId),
								  KEY user_client_idx (ClientId),
								  CONSTRAINT user_client FOREIGN KEY (ClientId) REFERENCES client (ClientId) ON DELETE NO ACTION ON UPDATE NO ACTION
							) ENGINE=InnoDB DEFAULT CHARSET=utf8;`)

	if err != nil {
		panic(err)
	}
}

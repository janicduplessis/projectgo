package interfaces

type DbInitializerRepo DbRepo

const (
	DbVersion    int  = 2    // Will DELETE all data when you increment this number
	CheckVersion bool = true // Set to false in production environnement to ignore DbVersions
)

func NewDbInitializerRepo(dbHandlers map[string]DbHandler) *DbInitializerRepo {
	dbHandler := dbHandlers["DbInitializerRepo"]

	return &DbInitializerRepo{
		dbHandlers: dbHandlers,
		dbHandler:  dbHandler,
	}
}

func (repo *DbInitializerRepo) Init() {

	// Check database version
	_, err := repo.dbHandler.Execute(`CREATE TABLE IF NOT EXISTS db_info (
										Version int(11) NOT NULL
									 );`)

	var version int
	err = repo.dbHandler.QueryRow(`SELECT * FROM db_info;`).Scan(&version)
	if err != nil {
		if err == ErrNoRows {
			// No db config yet
			_, err = repo.dbHandler.Execute(`INSERT INTO db_info VALUES(?);`, DbVersion)
			if err != nil {
				panic(err)
			}
			version = -1
		} else {
			panic(err)
		}
	}

	if version != DbVersion && CheckVersion {
		// Drop all tables
		_, err = repo.dbHandler.Execute(`DROP TABLE IF EXISTS client, user, channel, message;`)
		if err != nil {
			panic(err)
		}
		// Update the version
		_, err = repo.dbHandler.Execute(`UPDATE db_info SET Version = ?;`, DbVersion)
		if err != nil {
			panic(err)
		}
	}

	// Create the tables if it doesnt exists
	_, err = repo.dbHandler.Execute(`CREATE TABLE IF NOT EXISTS client (
										ClientId int(11) NOT NULL,
										DisplayName varchar(255) NOT NULL,
										FirstName varchar(255) NOT NULL,
										LastName varchar(255) NOT NULL,
										Email varchar(255) NOT NULL,
										PRIMARY KEY (ClientId)
									 );`)

	if err != nil {
		panic(err)
	}

	_, err = repo.dbHandler.Execute(`CREATE TABLE IF NOT EXISTS user (
								  		UserId int(11) NOT NULL AUTO_INCREMENT,
								 	 	Username varchar(64) NOT NULL,
								  		PasswordHash varchar(64) NOT NULL,
								  		PRIMARY KEY (UserId)
									 );`)

	if err != nil {
		panic(err)
	}

	_, err = repo.dbHandler.Execute(`CREATE TABLE IF NOT EXISTS channel (
								  	 	ChannelId int(11) NOT NULL AUTO_INCREMENT,
								  	 	Name varchar(255) NOT NULL,
								  	 	PRIMARY KEY (ChannelId)
									 );`)

	if err != nil {
		panic(err)
	}

	_, err = repo.dbHandler.Execute(`CREATE TABLE IF NOT EXISTS ct.message (
										MessageId INT NOT NULL AUTO_INCREMENT,
										Body TEXT NOT NULL,
										Time DATETIME NOT NULL,
										ClientId INT NOT NULL,
										ChannelId INT NOT NULL,
										PRIMARY KEY (MessageId),
										INDEX message_client_idx (ClientId ASC),
										INDEX message_channel_idx (ChannelId ASC),
										CONSTRAINT message_client
										    FOREIGN KEY (ClientId)
										    REFERENCES ct.client (ClientId)
										    ON DELETE NO ACTION
										    ON UPDATE NO ACTION,
										CONSTRAINT message_channel
										    FOREIGN KEY (ChannelId)
										    REFERENCES ct.channel (ChannelId)
										    ON DELETE NO ACTION
										    ON UPDATE NO ACTION
									);`)

	if err != nil {
		panic(err)
	}
}

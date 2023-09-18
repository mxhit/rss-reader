package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DRIVER_NAME = "sqlite3"
	DIR_PATH    = "./database/"
	DB_FILE     = "rss.db"
	FILE_RSS_DB = DIR_PATH + DB_FILE
	ERR_GENERIC = "Something went wrong: %s\n"
)

var sqlDb *sql.DB

func init() {
	sqlDb, _ = sql.Open(DRIVER_NAME, FILE_RSS_DB)
}

// Creates a `.db` file
func createDatabaseFile(filename string) {
	log.Printf("Creating %s\n", DB_FILE)

	if _, err := os.Create(filename); err != nil {
		log.Panicf("Error while creating %s\n: %v", DB_FILE, err.Error())
	}

	log.Printf("Successfully created %s\n", DB_FILE)
}

// checks if the file passed exists
func isExists(filename string) bool {
	if _, err := os.Open(filename); err != nil {
		return false
	}

	return true
}

func getAll(items map[string]string) map[string]string {
	feedEntities := make(map[string]string, 0)

	for author, _ := range items {
		row, err := sqlDb.Query("SELECT title FROM feed WHERE author = ?;", author)
		if err != nil {
			log.Panicf(ERR_GENERIC, err.Error())
		}

		defer row.Close()

		var title string
		for row.Next() && err == nil {
			row.Scan(&title)
		}

		feedEntities[author] = title
	}

	return feedEntities
}

func save(items map[string]string) {
	log.Println("Saving items to the database")

	insertFeedQuery := `INSERT INTO feed("author", "title") VALUES(?, ?);`

	statement, err := sqlDb.Prepare(insertFeedQuery)
	if err != nil {
		log.Panicf(ERR_GENERIC, err.Error())
	}

	for k, v := range items {
		if _, err := statement.Exec(k, v); err != nil {
			log.Panicf("Error when inserting records: %s\n", err.Error())
		}
	}
}

func createTable() {
	log.Println("Creating table 'feed' if it does not exist")

	createFeedQuery := `CREATE TABLE IF NOT EXISTS feed (id INTEGER PRIMARY KEY, author TEXT NOT NULL, title TEXT NOT NULL);`

	statement, err := sqlDb.Prepare(createFeedQuery)
	if err != nil {
		log.Panicf("Something went wrong while preparing the statement: %s\n", err.Error())
	}

	_, err = statement.Exec()
	if err != nil {
		log.Panicf("Something went wrong while creating table: %s\n", err.Error())
	}

	statement.Close()
}

func UpdateTable(toUpdate map[string]string) {
	log.Println("New entries available. Updating the database...")

	err := sqlDb.Ping()
	if err == nil {
		updateFeedQuery := `UPDATE feed SET title = ? WHERE author = ?;`

		statement, err := sqlDb.Prepare(updateFeedQuery)
		if err != nil {
			log.Panicf(ERR_GENERIC, err.Error())
		}

		for author, title := range toUpdate {
			if _, err := statement.Exec(title, author); err != nil {
				log.Panicf("Error when inserting records: %s\n", err.Error())
			}
		}
	} else {
		log.Panicf(ERR_GENERIC, err.Error())
	}
	log.Println("Updated the database")

	defer sqlDb.Close()
}

func GetExistingData(items map[string]string) map[string]string {
	if err := sqlDb.Ping(); err == nil {
		if !isExists(FILE_RSS_DB) {
			createDatabaseFile(FILE_RSS_DB)
			createTable()
			save(items)
		}

		lastEntries := getAll(items)

		return lastEntries
	} else {
		log.Panicf(ERR_GENERIC, err.Error())
	}

	defer sqlDb.Close()

	return nil
}

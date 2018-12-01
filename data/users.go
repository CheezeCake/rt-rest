package data

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/CheezeCake/rt-rest/config"
	_ "github.com/mattn/go-sqlite3"
)

const (
	usersTable = "users"
)

var (
	ErrUserNotFound      = errors.New("User not found")
	ErrUserAlreadyExists = errors.New("User already exists")
)

var (
	users Users
)

type Users struct {
	db              *sql.DB
	getPasswordStmt *sql.Stmt
	addUserStmt     *sql.Stmt
}

func AddUser(username, passwordHash string) error {
	_, err := GetPassword(username)
	if err == nil {
		return ErrUserAlreadyExists
	}
	if err != ErrUserNotFound {
		return err
	}
	_, err = users.addUserStmt.Exec(username, passwordHash)
	return err
}

func GetPassword(username string) (string, error) {
	var password string
	err := users.getPasswordStmt.QueryRow(username).Scan(&password)
	if err == sql.ErrNoRows {
		return password, ErrUserNotFound
	}
	return password, err
}

func CloseUsers() {
	users.getPasswordStmt.Close()
	users.db.Close()
}

func InitUsers() {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalln(err)
	}

	sql := fmt.Sprintf(
		"create table if not exists %s (id integer not null primary key, username text, password text);",
		usersTable)
	_, err = db.Exec(sql)
	if err != nil {
		db.Close()
		log.Fatalln(err)
	}

	sql = fmt.Sprintf("select password from %s where username = ?", usersTable)
	getPasswordStmt, err := db.Prepare(sql)
	if err != nil {
		db.Close()
		log.Fatalln(err)
	}

	sql = fmt.Sprintf("insert into %s (username, password) values(?, ?);", usersTable)
	addUserStmt, err := db.Prepare(sql)
	if err != nil {
		getPasswordStmt.Close()
		db.Close()
		log.Fatalln(err)
	}

	users = Users{db: db, getPasswordStmt: getPasswordStmt, addUserStmt: addUserStmt}
}

func Init(cfg config.Cfg) {
	InitUsers()
}

package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type userInfo struct {
	username string
	password string
	email    string
}

type roomInfo struct {
	roomName string
	public   bool
	roomPass string
}

var db *sql.DB

func dbConnect(dbUser, dbPass, dbName string) {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName))
	checkErr(err)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

//Check user's credentials
func checkCredentials(username, password string) bool {
	// query
	rows, err := db.Query("SELECT * FROM Users WHERE username=?", username)
	checkErr(err)

	user := userInfo{}

	for rows.Next() {
		err = rows.Scan(&user.username, &user.password, &user.email)
		checkErr(err)
		if user.username == username && user.password == md5Hash(password) {
			return true
		}
	}
	return false
}

//Get a room's info
func getRoom(roomName string) roomInfo {
	// query
	rows, err := db.Query("SELECT roomName,public,roomPass FROM Rooms WHERE roomName=?", roomName)
	checkErr(err)

	room := roomInfo{}
	room.roomName = ""
	for rows.Next() {
		err = rows.Scan(&room.roomName, &room.public, &room.roomPass)
		checkErr(err)
	}
	return room
}

//Create a new room
func createNewRoom(roomName, password string, public bool) {
	//insert
	stmt, err := db.Prepare("INSERT Rooms SET roomName=?,public=?,roomPass=?")
	checkErr(err)

	_, err = stmt.Exec(roomName, public, password)
	checkErr(err)
}

//Returns the md5Hash of string s
func md5Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

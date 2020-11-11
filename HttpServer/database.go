package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func dbConnect(dbUser, dbPass, dbName string) {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName))
	checkErr(err)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

//Get a user from the database
func getUser(username string) userInfo {
	// query
	rows, err := db.Query("SELECT * FROM Users WHERE username=?", username)
	checkErr(err)

	user := userInfo{}
	user.Username = ""

	for rows.Next() {
		err = rows.Scan(&user.Username, &user.password, &user.Email)
		checkErr(err)
		fmt.Println(user)
	}
	return user
}

//Returns the md5Hash of string s
func md5Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

//Add a user in the database
func addUser(user userInfo) {
	//Hash the password before inserting it in the database
	passHash := md5Hash(user.password)
	fmt.Println(passHash)
	//insert
	stmt, err := db.Prepare("INSERT Users SET username=?,password=?,email=?")
	checkErr(err)

	_, err = stmt.Exec(user.Username, passHash, user.Email)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

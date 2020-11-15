package main

import "net"

//Handle Login Request
func handleLogin(conn net.Conn, data []string) {
	username := data[1]
	password := data[2]

	if checkCredentials(username, password) {
		sendMsg(conn, "success\n")
	} else {
		sendMsg(conn, "failed\n")
	}
}

//Handle JoinRoom request
func handleJoinRoom(conn net.Conn, data []string) {
	roomName := data[1]
	room := getRoom(roomName)
	//If room doesn't exitst
	if room.roomName == "" {
		sendMsg(conn, "failed\n")
	} else { //Check if is a private room
		if room.public == false {
			sendMsg(conn, "password\n")
		} else {
			sendMsg(conn, "success\n")
		}
	}
}

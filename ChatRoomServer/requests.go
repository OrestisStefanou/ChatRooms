package main

import (
	"net"
)

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
	username := data[1]
	roomName := data[2]
	room := getRoom(roomName)
	//If room doesn't exitst
	if room.roomName == "" {
		sendMsg(conn, "failed\n")
	} else { //Check if is a private room
		if room.public == false {
			sendMsg(conn, "password\n")
			//Read the password,check if is correct
			password := recMsg(conn)
			if room.roomPass == password {
				sendMsg(conn, "success\n")
				enterRoom(conn, username, room.roomName)
			} else {
				sendMsg(conn, "Wrong Password\n")
			}
		} else {
			sendMsg(conn, "success\n")
			enterRoom(conn, username, room.roomName)
			//fmt.Printf("ChatRooms is %#v\n", chatRooms)
		}
	}
}

//Handle a client room message
func handleMessage(conn net.Conn, data []string) {
	userName := data[1]
	roomName := data[2]
	message := data[3]

	//Insert the messageInfo in the room's channel
	msg := userConn{message, conn, userName, roomName}
	//MUTEX HERE ??
	chatRoomMsgs[roomName] <- msg
}

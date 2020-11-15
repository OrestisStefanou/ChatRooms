package main

import (
	"fmt"
	"net"
)

func printMenu() {
	fmt.Println("1.Join a room")
	fmt.Println("2.Create a room")
	fmt.Println("3.Exit")
}

func joinRoom(conn net.Conn) {
	fmt.Print("Enter room name:")
	roomName := readString()
	request := fmt.Sprintf("JoinRoom%s%s\n", specialString, roomName)
	sendMsg(conn, request)
	response := recMsg(conn)
	switch response {
	case "success":
		fmt.Println("Entering room")
	case "password":
		//Ask for password and send it to the server
		fmt.Println("Asking for password")
	case "failed":
		//Room name doesn't exist
		fmt.Println("Room name doesn't exist")
	}
}

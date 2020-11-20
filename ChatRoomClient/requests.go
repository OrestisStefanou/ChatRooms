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

//Create a room
func createRoom(conn net.Conn) {
	fmt.Print("Enter Room Name:")
	roomName := readString()
	fmt.Print("Private(0) or Public(1)?Enter 0 or 1:")
	public := readString()
	switch public {
	case "1":
		//Send a request to create a public room
		request := fmt.Sprintf("CreateRoom%s%s%s%s%s''\n", specialString, roomName, specialString, public, specialString)
		sendMsg(conn, request)
		//Get the response
		response := recMsg(conn)
		fmt.Println(response)
	case "0":
		fmt.Print("Enter room's password:")
		password := readString()
		//Send a request to create a private room
		request := fmt.Sprintf("CreateRoom%s%s%s%s%s%s\n", specialString, roomName, specialString, public, specialString, password)
		sendMsg(conn, request)
		//Get the response
		response := recMsg(conn)
		fmt.Println(response)
	default:
		fmt.Print("Please enter 0 or 1")
	}
}

//Enter a room
func enterRoom(conn net.Conn) {
	ch := make(chan int, 1) //A channel to know when client exits the room
	//Start chatting
	go func() { //Go routine to get incoming messages
		for {
			select {
			case _ = <-ch:
				goto end
			default:
				incomingMsg := recMsg(conn)
				if incomingMsg == "Exit successful" {

				} else {
					fmt.Println(incomingMsg)
				}
			}
		}

	end:
		//fmt.Println("Go routine exits")
	}()
	msgRequest := fmt.Sprintf("Message%s%s%s%s%sJoined the room\n", specialString, myUsername, specialString, myRoom, specialString)
	sendMsg(conn, msgRequest)
	fmt.Print("Send msg(send 'exit' to exit room ):\n")
	for {
		message := readString()
		if message == "exit" {
			//Send the exit message to the server
			msgRequest := fmt.Sprintf("Message%s%s%s%s%s%s\n", specialString, myUsername, specialString, myRoom, specialString, message)
			sendMsg(conn, msgRequest)
			ch <- 1
			break
		} else {
			msgRequest := fmt.Sprintf("Message%s%s%s%s%s%s\n", specialString, myUsername, specialString, myRoom, specialString, message)
			sendMsg(conn, msgRequest)
		}
	}
}

func joinRoom(conn net.Conn) {
	fmt.Print("Enter room name:")
	roomName := readString()
	request := fmt.Sprintf("JoinRoom%s%s%s%s\n", specialString, myUsername, specialString, roomName)
	sendMsg(conn, request)
	response := recMsg(conn)
	switch response {
	case "success":
		fmt.Println("Entering room")
		myRoom = roomName
		enterRoom(conn)
	case "password":
		//Ask for password and send it to the server
		fmt.Print("Room Password:")
		roomPass := readString()
		request := fmt.Sprintf("%s\n", roomPass)
		sendMsg(conn, request)
		resp := recMsg(conn)
		if resp == "success" {
			fmt.Println("Entering room")
			myRoom = roomName
			enterRoom(conn)
		} else {
			fmt.Println(resp)
		}
	case "failed":
		//Room name doesn't exist
		fmt.Println("Room name doesn't exist")
	}
}

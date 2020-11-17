package main

import (
	"fmt"
	"net"
)

type msgInfo struct {
	msg        string
	sender     net.Conn
	senderName string
}

//A chatRoom is a set of client connections
var chatRooms map[string][]net.Conn

//A map with a chanel of msg for each room
var chatRoomMsgs map[string]chan msgInfo

//Initialize the maps
func initChatRooms() {
	chatRooms = make(map[string][]net.Conn)
	chatRoomMsgs = make(map[string]chan msgInfo)
}

//Create a new room
func createRoom(roomName string) {
	chatRooms[roomName] = make([]net.Conn, 0)
	chatRoomMsgs[roomName] = make(chan msgInfo)
}

//Add a user in the room
func addUserInRoom(roomName string, clientConn net.Conn) {
	chatRooms[roomName] = append(chatRooms[roomName], clientConn)
}

//Handle a Room
func handleRoom(roomName string) {
	var channel chan msgInfo = chatRoomMsgs[roomName]
	var slice []net.Conn = chatRooms[roomName]
	var msg string
	for {
		select {
		case m := <-channel: //In case there is a msg in the channel
			slice = chatRooms[roomName]
			fmt.Printf("Received msg:%s from:", m.msg)
			fmt.Println(m.sender)
			if m.msg == "exit" {
				msg = fmt.Sprintf("[%s] left the room\n", m.senderName)
				//Find a way to remove clients who exited()
			} else {
				msg = fmt.Sprintf("[%s]:%s\n", m.senderName, m.msg)
			}
			fmt.Printf("Slice is %#v\n", slice)
			//send it to the others users in the room
			for _, user := range slice {
				if user != m.sender {
					sendMsg(user, msg)
				}
			}
		default:
			//fmt.Println("Checking if the room is empty")
			if len(slice) == 0 {
				fmt.Println("Exiting the chanel")
				goto end
			}
		}
	}
end:
}

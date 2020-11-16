package main

import "net"

//A chatRoom is a set of client connections
var chatRooms map[string][]net.Conn

//A map with the messages to send to each room
var chatRoomMsgs map[string]chan string

//Initialize the maps
func initChatRooms() {
	chatRooms = make(map[string][]net.Conn)
	chatRoomMsgs = make(map[string]chan string)
}

//Create a new room
func createRoom(roomName string) {
	chatRooms[roomName] = make([]net.Conn, 0)
}

//Add a user in the room
func addUserInRoom(roomName string, clientConn net.Conn) {
	chatRooms[roomName] = append(chatRooms[roomName], clientConn)
}

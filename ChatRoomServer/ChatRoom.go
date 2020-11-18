package main

import (
	"fmt"
	"net"
	"sync"
)

type userConn struct {
	msg      string
	conn     net.Conn
	userName string
	room     string
}

//A chatRoom is a set of client connections
var chatRooms map[string][]userConn

//A map with a chanel of msgs for each room
var chatRoomMsgs map[string]chan userConn

//A mutex to protect ChatRooms map
var mu sync.Mutex

//Initialize the maps
func initChatRooms() {
	chatRooms = make(map[string][]userConn)
	chatRoomMsgs = make(map[string]chan userConn)
}

//Create a new room
func createRoom(roomName string) {
	chatRooms[roomName] = make([]userConn, 0)
	chatRoomMsgs[roomName] = make(chan userConn)
}

//Add a user in the room
func addUserInRoom(roomName string, clientConn userConn) {
	mu.Lock()
	chatRooms[roomName] = append(chatRooms[roomName], clientConn)
	mu.Unlock()
}

//User enters the room
func enterRoom(conn net.Conn, username, roomName string) {
	userConnection := userConn{"", conn, username, roomName}
	//Check if the room exists
	_, hasKey := chatRooms[roomName]
	if hasKey { //Append the user in the room
		addUserInRoom(roomName, userConnection)
		//fmt.Printf("Slice is %#v\n", chatRooms[room.roomName])

	} else { //Create the room and append the user
		createRoom(roomName)
		addUserInRoom(roomName, userConnection)
		//Start a go routine to handle the room
		go handleRoom(roomName)
	}
}

//Handle a Room
func handleRoom(roomName string) {
	var msgs chan userConn = chatRoomMsgs[roomName]
	var userConns []userConn = chatRooms[roomName]
	var msg string
	for {
		select {
		case m := <-msgs: //In case there is a msg in the channel
			userConns = chatRooms[roomName] //In case new users joined the room
			fmt.Printf("Received msg:%s from:", m.msg)
			fmt.Println(m.userName)
			if m.msg == "exit" {
				msg = fmt.Sprintf("[%s] left the room\n", m.userName)
				sendMsg(m.conn, "Exit successful\n")
				//Remove the user from the group
				for i := 0; i < len(userConns); i++ {
					if userConns[i].userName == m.userName {
						userConns[i].room = ""
						break
					}
				}
			} else {
				msg = fmt.Sprintf("[%s]:%s\n", m.userName, m.msg)
			}
			//fmt.Printf("userConns is %#v\n", userConns)
			//send the message to the other users in the room
			for _, user := range userConns {
				if user.conn != m.conn && user.room == roomName {
					sendMsg(user.conn, msg)
				}
			}
		default:
			//Check if room is empty
			exitFlag := true
			for i := 0; i < len(userConns); i++ {
				if userConns[i].room != "" {
					exitFlag = false
				}
			}
			if exitFlag {
				fmt.Println("Exiting the chanel")
				close(chatRoomMsgs[roomName])
				delete(chatRooms, roomName)
				delete(chatRoomMsgs, roomName)
				goto end
			}
		}
	}
end:
}

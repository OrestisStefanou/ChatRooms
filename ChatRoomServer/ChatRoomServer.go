package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", listeningPort)
	checkError(err)
	//Listen for connections
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	dbConnect(dbUser, dbPass, dbName) //Connect to the database

	//Main loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	finished := false
	for finished == false {
		request := recMsg(conn)
		data := strings.Split(request, specialString)
		fmt.Printf("data: %#v\n", data)
		switch data[0] {
		case "Login":
			handleLogin(conn, data)
		case "Closed": //Connection finished
			finished = true
		default:
			sendMsg(conn, "Something went wrong\n")
		}
	}
	//sendMsg(conn, "Hello from the server!\n")
	// we're finished with this client
}

//Handle Login Request
func handleLogin(conn net.Conn, data []string) {
	username := data[1]
	password := data[2]
	fmt.Printf("USERNAME:%s\n", username)
	fmt.Printf("PASS:%s\n", password)
	if checkCredentials(username, password) {
		sendMsg(conn, "success\n")
	} else {
		sendMsg(conn, "failed\n")
	}
}

//Send msg to server
func sendMsg(conn net.Conn, msg string) {
	bytesSent := 0
	msgLen, err := fmt.Fprintf(conn, msg)
	checkError(err)
	bytesSent = bytesSent + msgLen
	for bytesSent != len(msg) {
		msgLen, err = fmt.Fprintf(conn, msg[bytesSent:])
		checkError(err)
		bytesSent = bytesSent + msgLen
	}
}

//Receive msg froms server
func recMsg(conn net.Conn) string {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Sprintf("Closed%s", specialString)
	}
	return strings.TrimSuffix(message, "\n")
}

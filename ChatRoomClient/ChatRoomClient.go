package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	//Get this from conf file
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage:%s host:port", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	//Connect to the server
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	//Login
	msg := "Hello from the client!\n"
	sendMsg(conn, msg)
	//fmt.Fprintf(conn, "Hello from the client!\n")
	//Receive msg from server
	message := recMsg(conn)
	checkError(err)
	fmt.Println(message)
	username, password := credentials()
	fmt.Printf("\nUsername: %s, Password: %s\n", username, password)
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

//Get the credentials of a user
func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(0)
	checkError(err)
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
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
	checkError(err)
	return message
}

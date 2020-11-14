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
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverService)
	checkError(err)
	//Connect to the server
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	//Login
	loggedIn := false
	for loggedIn == false {
		loggedIn = login(conn)
	}
	os.Exit(0)
}

func login(conn net.Conn) bool {
	username, password := credentials() //Get user's credential from user
	msg := fmt.Sprintf("Login%s%s%s%s\n", specialString, username, specialString, password)
	sendMsg(conn, msg) //Send the request
	response := recMsg(conn)
	fmt.Printf("\n%s\n", response)
	if response == "success" {
		return true
	}
	return false
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
	return strings.TrimSuffix(message, "\n")
}

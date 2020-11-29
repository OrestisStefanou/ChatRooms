package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var myUsername = ""
var myRoom = ""

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

	timeToExit := false
	for timeToExit == false {
		printMenu()
		fmt.Print("Enter a choice(1-3):")
		//Get user's request
		text := readString()
		switch text {
		case "1":
			joinRoom(conn)
		case "2":
			createRoom(conn)
		case "3":
			//Exit
			timeToExit = true
		default:
			fmt.Println("Please enter a number between 1-3")
		}
	}
	os.Exit(0)
}

func readString() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func login(conn net.Conn) bool {
	username, password := credentials() //Get user's credential from user
	msg := fmt.Sprintf("Login%s%s%s%s\n", specialString, username, specialString, password)
	sendMsg(conn, msg) //Send the request
	response := recMsg(conn)
	//fmt.Printf("\n%s\n", response)
	if response == "success" {
		fmt.Println("Login successful")
		myUsername = username
		return true
	}
	fmt.Println("Login failed")
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

func cipher(text string, direction int) string {
	// shift -> number of letters to move to right or left
	// offset -> size of the alphabet, in this case the plain ASCII
	shift, offset := rune(3), rune(26)

	// string->rune conversion
	runes := []rune(text)

	for index, char := range runes {
		// Iterate over all runes, and perform substitution
		// wherever possible. If the letter is not in the range
		// [1 .. 25], the offset defined above is added or
		// subtracted.
		switch direction {
		case -1: // encoding
			if char >= 'a'+shift && char <= 'z' ||
				char >= 'A'+shift && char <= 'Z' {
				char = char - shift
			} else if char >= 'a' && char < 'a'+shift ||
				char >= 'A' && char < 'A'+shift {
				char = char - shift + offset
			}
		case +1: // decoding
			if char >= 'a' && char <= 'z'-shift ||
				char >= 'A' && char <= 'Z'-shift {
				char = char + shift
			} else if char > 'z'-shift && char <= 'z' ||
				char > 'Z'-shift && char <= 'Z' {
				char = char + shift - offset
			}
		}

		// Above `if`s handle both upper and lower case ASCII
		// characters; anything else is returned as is (includes
		// numbers, punctuation and space).
		runes[index] = char
	}

	return string(runes)
}

// encode and decode provide the API for encoding and decoding text using
// the Caesar Cipher algorithm.
func encode(text string) string { return cipher(text, -1) }
func decode(text string) string { return cipher(text, +1) }

//Send msg
func sendMsg(conn net.Conn, msg string) {
	msgToSend := encode(msg) + "\n"
	bytesSent := 0
	msgLen, err := fmt.Fprintf(conn, msgToSend)
	checkError(err)
	bytesSent = bytesSent + msgLen
	for bytesSent != len(msgToSend) {
		msgLen, err = fmt.Fprintf(conn, msgToSend[bytesSent:])
		checkError(err)
		bytesSent = bytesSent + msgLen
	}
}

//Receive msg
func recMsg(conn net.Conn) string {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Sprintf("Closed%s", specialString)
	}
	return strings.Trim(decode(message), "\n")
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

//String to split the request data
var specialString = "@*@"

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
	bytesSent = msgLen
	for bytesSent < len(msgToSend) {
		msgLen, err = fmt.Fprintf(conn, msgToSend[bytesSent:])
		checkError(err)
		bytesSent = bytesSent + msgLen
	}
}

//Receive msg
func recMsg(conn net.Conn) string {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Sprintf("Closed")
	}
	return strings.Trim(decode(message), "\n")
}

func printYellow(text string) {
	colorYellow := "\033[33m"
	fmt.Println(string(colorYellow), text)
	colorReset := "\033[0m"
	fmt.Println(string(colorReset))
}

func printGreen(text string) {
	colorGreen := "\033[32m"
	fmt.Println(string(colorGreen), text)
	colorReset := "\033[0m"
	fmt.Println(string(colorReset))
}

func printRed(text string) {
	colorRed := "\033[31m"
	fmt.Println(string(colorRed), text)
	colorReset := "\033[0m"
	fmt.Println(string(colorReset))
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var ch chan string //Channel to send messages to goroutine that handles the log file

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", listeningPort)
	checkError(err)
	//Listen for connections
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	go handleLogFile()
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
	var b bytes.Buffer
	for finished == false {
		message := recMsg(conn)
		//fmt.Printf("Got message:%s\n", message)
		data := strings.Split(message, specialString)
		switch data[0] {
		case "MemInfo":
			//Process and print MemInfo
			//fmt.Println(data[1])
			b.WriteString(data[1] + "\n")
			sendMsg(conn, "Got it\n")
		case "CpuInfo":
			//Process and print CpuInfo
			//fmt.Println(data[1])
			b.WriteString(data[1] + "\n")
			sendMsg(conn, "Got it\n")
		case "ClientsNum":
			//Print how many clients are connected to the server
			//fmt.Println(data[1])
			b.WriteString(data[1] + "\n")
			sendMsg(conn, "Got it\n")
		case "StatsDone":
			printYellow(data[1])
			fmt.Println(b.String())
			ch <- data[1] + b.String()
			b.Reset()
			sendMsg(conn, "Got it\n")
		default:
			finished = true
		}
	}
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

func handleLogFile() {
	ch = make(chan string)
	currentTime := time.Now()
	//Open or create the log file
	path := logDir + currentTime.Format("01-02-2006")
	logfile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("os.OpenFile() failed with '%s\n", err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	for {
		stat := <-ch
		log.Println(stat)
	}

}

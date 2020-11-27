package main

import (
	"bufio"
	"bytes"
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
			b.Reset()
			sendMsg(conn, "Got it\n")
		default:
			finished = true
		}
	}
}

//Send msg
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

//Receive msg
func recMsg(conn net.Conn) string {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Sprintf("Closed%s", specialString)
	}
	return strings.TrimSuffix(message, "\n")
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

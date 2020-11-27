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
		message := recMsg(conn)
		//fmt.Printf("Got message:%s\n", message)
		data := strings.Split(message, specialString)
		switch data[0] {
		case "MemInfo":
			//Process and print MemInfo
			fmt.Println(data[1])
			sendMsg(conn, "Got it\n")
		case "CpuInfo":
			//Process and print CpuInfo
			fmt.Println(data[1])
			sendMsg(conn, "Got it\n")
		case "ClientsNum":
			//Print how many clients are connected to the server
			fmt.Println(data[1])
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

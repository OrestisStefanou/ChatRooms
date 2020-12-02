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
			b.WriteString(data[1] + "\n")
			sendMsg(conn, "Got it\n")
		case "CpuInfo":
			b.WriteString(data[1] + "\n")
			sendMsg(conn, "Got it\n")
		case "ClientsNum":
			b.WriteString(data[1] + "\n")
			sendMsg(conn, "Got it\n")
		case "StatsDone":
			printYellow(data[1])
			fmt.Println(b.String())
			ch <- data[1] + b.String()
			b.Reset()
			sendMsg(conn, "Got it\n")
		case "GetUpdates":
			handleUpdates(conn)
		case "CheckForUpdates":
			msg := fmt.Sprintf("%s\n", myVersion)
			sendMsg(conn, msg)
		default:
			finished = true
		}
	}
}

//I am uising a different specialString here because the original exists
//in the files that we send
func handleUpdates(conn net.Conn) {
	for i := 0; i < len(sharedFiles); i++ {
		filename := sharedFiles[i]
		//First send the filename to create
		msg := fmt.Sprintf("CreateFile%s%s\n", "&*&", filename)
		sendMsg(conn, msg)
		recMsg(conn) //To make sure the client got the message
		//Send the file line by line
		sendFile(filename, conn)
		//Send message that file is sent
		msg = fmt.Sprintf("CloseFile%s\n", "&*&")
		sendMsg(conn, msg)
		recMsg(conn)
	}
	//Send message that updates are done
	msg := fmt.Sprintf("Finished%s\n", "&*&")
	sendMsg(conn, msg)
	recMsg(conn)

}

func sendFile(filePath string, conn net.Conn) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		sendMsg(conn, "Something went wrong\n")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		msg := fmt.Sprintf("Line%s%s\n", "&*&", line)
		sendMsg(conn, msg) //Send the line
		recMsg(conn)       //Make sure the client got the message
	}
	if err = scanner.Err(); err != nil {
		sendMsg(conn, "Something went wrong\n")
	}
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

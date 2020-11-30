package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	getUpdates("./")
}

func getUpdates(dir string) {
	//Connect to the systemMonitor Server
	sysMonitorAddr, err := net.ResolveTCPAddr("tcp4", systemMonitorService)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, sysMonitorAddr)
	checkError(err)

	var f *os.File
	//Send a request to get the updates
	msg := fmt.Sprintf("GetUpdates%s\n", specialString)
	sendMsg(conn, msg)
	response := recMsg(conn)
	for {
		data := strings.Split(response, specialString)
		switch data[0] {
		case "CreateFile":
			filename := data[1]
			path := filepath.Join(dir, filename)
			//Create a file
			f, err = os.Create(path)
			checkError(err)

			sendMsg(conn, "Got it\n") //Tell the server that we got the message
			response = recMsg(conn)   //Receive a new message
		case "Line":
			//Append a line to the file
			f.WriteString(data[1] + "\n")
			sendMsg(conn, "Got it\n") //Tell the server that we got the message
			response = recMsg(conn)   //Receive a new message
		case "CloseFile":
			//Close the file
			f.Close()
			sendMsg(conn, "Got it\n") //Tell the server that we got the message
			response = recMsg(conn)   //Receive a new message
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {
	wg.Add(len(filesDir))
	for i := 0; i < len(filesDir); i++ {
		go getUpdates(filesDir[i])
	}
	//Wait for the goroutines to finish
	wg.Wait()
}

func getUpdates(dir string) {
	defer wg.Done()
	fmt.Println("Entering the function")
	//Connect to the systemMonitor Server
	sysMonitorAddr, err := net.ResolveTCPAddr("tcp4", systemMonitorService)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, sysMonitorAddr)
	checkError(err)

	var f *os.File
	finished := false
	//Send a request to get the updates
	msg := fmt.Sprintf("GetUpdates%s\n", specialString)
	sendMsg(conn, msg)
	response := recMsg(conn)
	for finished == false {
		data := strings.Split(response, "&*&")
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
		case "Finished":
			finished = true
			sendMsg(conn, "Got it\n")
		default:
			finished = true
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

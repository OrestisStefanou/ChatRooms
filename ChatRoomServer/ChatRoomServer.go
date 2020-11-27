package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", listeningPort)
	checkError(err)
	//Listen for connections
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	dbConnect(dbUser, dbPass, dbName) //Connect to the database
	initChatRooms()

	go sendStats() //Start the goroutine to send the stats to system monitor server

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
		request := recMsg(conn)
		data := strings.Split(request, specialString)
		switch data[0] {
		case "Login":
			handleLogin(conn, data)
		case "JoinRoom":
			handleJoinRoom(conn, data)
		case "Message":
			handleMessage(conn, data)
		case "CreateRoom":
			handleCreateRoom(conn, data)
		case "Closed": //Connection finished
			finished = true
			//Decrease the number of connected Users
			usersMutex.Lock()
			connectedUsers--
			usersMutex.Unlock()
		default:
			sendMsg(conn, "Something went wrong\n")
		}
	}
	// we're finished with this client
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

//Send cpu,memory stats and user count to systemMonitor Server
func sendStats() {
	//Connect to the systemMonitor Server first
	sysMonitorAddr, err := net.ResolveTCPAddr("tcp4", systemMonitorService)
	checkErr(err)

	conn, err := net.DialTCP("tcp", nil, sysMonitorAddr)
	checkErr(err)
	for {
		//Execute the commands and send the results
		cmd2 := exec.Command("mpstat", "2", "2")
		out, err := cmd2.CombinedOutput()
		if err != nil {
			log.Fatalf("cmd.CombinedOutput() failed with '%s'\n", err)
		}
		//fmt.Printf("Output:\n%s\n", string(out))
		tempInfo := strings.Replace(string(out), "%", " ", -1) //Replace '%' because it causes problems to fprintf
		info := strings.Split(tempInfo, "\n")
		//Send the info to the system monitor server
		for i := 0; i < len(info); i++ {
			msg := fmt.Sprintf("CpuInfo%s%s%s%s\n", specialString, info[i], specialString, myName)
			sendMsg(conn, msg)
			recMsg(conn)

		}

		cmd := "cat /proc/meminfo | grep 'Mem'"
		out, err = exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatalf("Command failed with '%s'\n", err)
		}
		//fmt.Println(string(out))
		info = strings.Split(string(out), "\n")
		for i := 0; i < len(info); i++ {
			msg := fmt.Sprintf("MemInfo%s%s%s%s\n", specialString, info[i], specialString, myName)
			sendMsg(conn, msg)
			recMsg(conn)
		}

		//Send the connected users
		//We don't lock the mutex here because is not a problem if we read an older number
		msg := fmt.Sprintf("ClientsNum%sConnected Users:%d%s%s\n", specialString, connectedUsers, specialString, myName)
		sendMsg(conn, msg)
		recMsg(conn)

		//Send message to tell system monitor that stats are finished
		msg = fmt.Sprintf("StatsDone%s%s\n", specialString, myName)
		sendMsg(conn, msg)
		recMsg(conn)
	}
}

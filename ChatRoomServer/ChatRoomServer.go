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
	//Add the necessary info to connect in conf.go file
	cmd := "cat /proc/meminfo | grep 'Mem'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatalf("Command failed with '%s'\n", err)
	}
	fmt.Println(string(out))

	cmd2 := exec.Command("mpstat", "2", "5")
	out, err = cmd2.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.CombinedOutput() failed with '%s'\n", err)
	}
	fmt.Printf("Output:\n%s\n", string(out))

}

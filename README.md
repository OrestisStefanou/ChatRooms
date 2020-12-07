# ChatRooms
A terminal chat rooms application in golang

HttpServer dependencies:
    <br />-Mysql database
    <br />-github.com/go-sql-driver/mysql library


ChatRoomServer dependencies:
    <br />-Mysql database
    <br />-github.com/go-sql-driver/mysql library

ChatRoomClient dependencies:
    <br />-golang.org/x/crypto/ssh/terminal library

How to setup HtppServer and ChatRoomServer:
    <br />-Create a user with read and write privilages in mysql
    <br />-Create a database with name ChatRooms
    <br />-Open ChatRooms/HttpServer/conf.go and ChatRooms/ChatRoomServer/conf.go and update dbUser and dbPass
     with the creadentials of the user you created.
    <br />In the file ChatRooms/ChatRoomServer/conf.go update the systemMonitorService to the ip adress of the machine that is running the systemMonitor server
    <br />Do the same in the files ChatRooms/Updater/conf.go and ChatRooms/ChatRoomClient/conf.go

How to build:
    <br />-Run the command go build in the following directories:
        <br />-ChatRooms/ChatRoomClient/
        <br />-ChatRooms/ChatRoomServer/
        <br />-ChatRooms/HttpServer/
	<br />-ChatRooms/SystemMonitor/
	<br />-ChatRooms/Updater/

How to run and use it:
    <br />-To run the ChatRoomServer and ChatRoomClient in different machines open the file
      ChatRooms/ChatRoomClient/conf.go and change serverService to "<ip address of the machine that the server runs>:1200"
    <br />-First run the SystemMonitorServer using the command -> ./ChatRooms/SystemMonitor/SystemMonitor
    <br />-Run the httpServer using this command -> ./ChatRooms/HttpServer/HttpServer
    <br />-Open a browser and go to the address:<Ip address of the machince that http server runs>:8080
    <br />-Create an account
    <br />-Run the ChatRoomServer using this command -> ./ChatRooms/ChatRoomServer/ChatRoomServer
    <br />-Run the ChatRoomClient using this command -> ./ChatRooms/ChatRoomClient/ChatRoomClient
    <br />-Login,create or join a room and start chatting!

What is the use of Updater:
<br />Because the file communications.go is the same for the client and the server if we want to change something we edit the communications.go file
in the SystemMonitor and then you run the Updater to update the communications file of ChatRoomClient and ChatRoomServer.

What left:
    <br />1.Better design of the html files
    <br />2.Create a load balancer

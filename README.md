# ChatRooms
A terminal chat rooms application in golang

HttpServer dependencies:
    -Mysql database
    -github.com/go-sql-driver/mysql library


ChatRoomServer dependencies:
    -Mysql database
    -github.com/go-sql-driver/mysql library

ChatRoomClient dependencies:
    -golang.org/x/crypto/ssh/terminal library

How to setup HtppServer and ChatRoomServer:
    -Create a user with read and write privilages in mysql
    -Create a database with name ChatRooms
    -Open ChatRooms/HttpServer/conf.go and Chatrooms/ChatRoom/Server.go and update dbUser and dbPass
     with the creadentials of the user you created

How to build:
    -Run the command go build in the following directories:
        -ChatRooms/ChatRoomClient/
        -ChatRooms/ChatRoomServer/
        -ChatRooms/HttpServer/

How to run and use it:
    -To run the ChatRoomServer and ChatRoomClient in different machines open the file
      ChatRooms/ChatRoomClient/conf.go and change serverService to "<ip address of the machine that the server runs>:1200"
    1.Run the httpServer using this command -> ./ChatRooms/HttpServer/HttpServer
    2.Open a browser and go to the address:<Ip address of the machince that http server runs>:8080
    3.Create an account
    4.Run the ChatRoomServer using this command -> ./ChatRooms/ChatRoomServer/ChatRoomServer
    5.Run the ChatRoomClient using this command -> ./ChatRooms/ChatRoomClient/ChatRoomClient
    6.Login,create or join a room and start chatting!

What left:
    1.Better design of the html files
    2.Create a load balancer
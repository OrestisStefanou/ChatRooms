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
    <br />-Open ChatRooms/HttpServer/conf.go and Chatrooms/ChatRoom/Server.go and update dbUser and dbPass
     with the creadentials of the user you created

How to build:
    <br />-Run the command go build in the following directories:
        <br />-ChatRooms/ChatRoomClient/
        <br />-ChatRooms/ChatRoomServer/
        <br />-ChatRooms/HttpServer/

How to run and use it:
    <br />-To run the ChatRoomServer and ChatRoomClient in different machines open the file
      ChatRooms/ChatRoomClient/conf.go and change serverService to "<ip address of the machine that the server runs>:1200"
    <br />1.Run the httpServer using this command -> ./ChatRooms/HttpServer/HttpServer
    <br />2.Open a browser and go to the address:<Ip address of the machince that http server runs>:8080
    <br />3.Create an account
    <br />4.Run the ChatRoomServer using this command -> ./ChatRooms/ChatRoomServer/ChatRoomServer
    <br />5.Run the ChatRoomClient using this command -> ./ChatRooms/ChatRoomClient/ChatRoomClient
    <br />6.Login,create or join a room and start chatting!

What left:
    <br />1.Better design of the html files
    <br />2.Create a load balancer
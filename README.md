# chat_app

Client Server Chat Application in Golang 

<br>To start the application you need to use the command : ```go build .```
<br>Then you need separate terminals for the server and each client
<br>You can appropriately run it using ```./chat_app s``` (s - for server) or ```./chat_app c``` (c - for client)

<br>client can do the following functions : private messaging, broadcasting message and quitting the server
<br>private message another client : ```/pm {receiver name} {message content}```
<br>braodcast message to all others : ```/broad {message content}```
<br>leave the server : ```/quit```

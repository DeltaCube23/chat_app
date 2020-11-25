# chat_app

Client Server Chat Application in Golang 

<br>To start the application you need to use the command : go build .
<br>Then you need separate terminals for the server and each client
<br>You can appropriately run it using ./chat_app s/c (s - for server, c - for client)
<br> Eg: ./chat_app s (will start server)

<br>client can do the following functions and all the messages need to be terminated with $
<br>private message another client : /pm {receiver name} {message content}$
<br>braodcast message to all others : /broad {message content}$
<br>leave the server : /quit$

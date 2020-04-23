package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object
type Message struct {
	Username  string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	crontab := cron.NewWithSeconds()//不支持秒
	_, _ = crontab.AddFunc( "*/1800 * * * * ?",func() {
		clients = make(map[*websocket.Conn]bool)
		fmt.Println(1)
	})
	crontab.Start()

	// Create a simple file server
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Configure websocket route
	http.HandleFunc("/ws/v2", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(2, w , r)
	})

	http.HandleFunc("/ws/v25", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(25, w , r)
	})

	http.HandleFunc("/ws/vn", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(0, w , r)
	})

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(i int,w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	if i == 2 && len(clients) > 1 {
		_ = ws.WriteJSON("超出人数")
		return
	}else if i == 25 && len(clients) > 24{
		_ = ws.WriteJSON("err")
		return
	}

	// Register our new client
	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				_ = client.Close()
				delete(clients, client)
			}
		}
	}
}
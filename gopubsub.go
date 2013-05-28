package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	id       int
	dataChan chan string
	complete bool
}

func registerClient(newClient *Client, clientCounter *int, clients map[int]*Client) {

	if newClient.id == 0 {
		*clientCounter++
		newClient.id = *clientCounter
		clients[*clientCounter] = newClient
		fmt.Println("new client connected" + " " + fmt.Sprintf("%v", newClient.id))
	} else if newClient.complete == true {
		fmt.Println("existing client is finished" + " " + fmt.Sprintf("%v", newClient.id))
		delete(clients, newClient.id)
	}

}

func generateData(sendChan chan string) {
	counter := 0

	for {
		counter++
		sendChan <- "Ping" + fmt.Sprintf("%v", counter) + ":"
		time.Sleep(time.Second)
	}
}

func sendData(msg string, clients *map[int]*Client) {
	for _, client := range *clients {
		client.dataChan <- msg + fmt.Sprintf("%v", client.id)
	}
}

func clientRegister(clientsChan chan *Client, sendChan chan string) {

	clients := make(map[int]*Client)
	clientCounter := 0
	fmt.Println("Client register started...")

	for {
		select {
		//New client connects
		case newClient := <-clientsChan:
			registerClient(newClient, &clientCounter, clients)
		//New message to be sent to clients
		case msg := <-sendChan:
			sendData(msg, &clients)
		}
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request, clientsChan chan *Client) {
	w.Header().Set("Content-Type", "application/json")

	//Create new client object and send to register via chan.
	var myDataChan chan string = make(chan string)
	newClient := &Client{id: 0, dataChan: myDataChan}
	clientsChan <- newClient

	//Read from data chan any new messages.
	msg := <-newClient.dataChan
	data, _ := json.Marshal(msg)

	fmt.Fprintln(w, string(data))

	newClient.complete = true
	clientsChan <- newClient
}

func main() {

	var clientsChan chan *Client = make(chan *Client)
	var sendChan chan string = make(chan string, 100)

	go clientRegister(clientsChan, sendChan)
	go generateData(sendChan)

	http.Handle("/", http.FileServer(http.Dir("/var/www/html/gopubsub")))

	http.HandleFunc("/poll", func(w http.ResponseWriter, r *http.Request) {
		httpHandler(w, r, clientsChan)
	})
	http.ListenAndServe(":4000", nil)

}

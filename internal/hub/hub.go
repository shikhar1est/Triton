package hub

import "log"

type Hub struct { //Hub is a central message broker
	// for managing client connections and broadcasting messages.
	//It's basically the brain of a real-time communication system.
	//It  manages all connected clients,registers new clients,
	// unregisters disconnected clients,broadcasts messages to all clients.

	clients map[*Client]bool //Go has no built-in set type,
	//  so we use a map with bool values to represent a set of clients who are connected.
	//Here, clients is a map where the keys are pointers to Client structs,
	// and the values are booleans indicating whether the client is connected (true) or not (false).
	broadcast  chan *Message //message to be sent to all clients
	register   chan *Client  //channel for registering new clients
	unregister chan *Client  //channel for unregistering disconnected clients
}

func NewHub() *Hub { //returns a pointer to a new Hub instance
	return &Hub{ //this initializes the Hub struct with empty maps and channels.
		//make is a built-in function in Go that allocates and initializes objects like slices, maps, and channels.
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register: //when a new client registers
			h.clients[client] = true //add the client to the clients map
			log.Println("New Client Registered, Total client count:",len(h.clients))

		case clien:=<-h.unregister: //when a client unregisters
		if _,ok:=h.clients[client];ok{ //check if the client exists in the map	
		    delete(h.clients,client) //remove the client from the map
			close(client.send) //close the client's send channel
			log.Println("Client Unregistered, Total client count:",len(h.clients))
		}

	case message:=<-h.broadcast:
		log.Printf("Broadcasting message to %d clients, payload: %s\n",len(h.clients), message.Payload)
		for client:= range h.clients{ //iterate over all connected clients
			select{
				case client.send<-message: //try to send the message to the client's send channel
			default:
				log.Println("Client queue full, closing connection")
				close(client.send) //close the client's send channel
				delete(h.clients,client) //remove the client from the map
			}
		}
	}
	}
}

//in Go, map syntax is used to create a map data structure,
//like map[keyType]valueType, so for example,map[string]int is a map
// where the keys are strings and the values are integers.

//Channels are a way for goroutines to communicate with each other and synchronize their execution.
// "chan" is the keyword used to declare a channel type in Go.
//A channel lets you send and receive values of a specific type between goroutines.
//Example, var ch chan int declares a channel that can send and receive integers.
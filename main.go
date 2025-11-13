package main

import (
	"log"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter() //here we use the mux package to create a new router
	r.HandleFunc("/ws",handleWebSocket) //here we define a route for WebSocket connections
	//we are basically telling the router to listen for requests at the "/ws" endpoint and 
	//to handle those requests using the handleWebSocket function
    log.Println("Server starting on port 8000")
	if err:=r.ListenAndServe(":8000",nil);err!=nil{
		log.Fatal("ListenAndServe: ",err)
	}

}
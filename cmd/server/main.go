package main

import (
	"github.com/larking7days/websocket_chatbot/api"
	"github.com/larking7days/websocket_chatbot/database"
	"log"
	"net/http"
)

func main() {
	db := database.InitDB()
	defer db.Close()

	r := api.SetupRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}

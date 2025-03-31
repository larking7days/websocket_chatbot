package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/larking7days/websocket_chatbot/database"
	"github.com/larking7days/websocket_chatbot/internal/models"
	"github.com/larking7days/websocket_chatbot/internal/websocket"
	"net/http"
	"strconv"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// WebSocket端点
	r.HandleFunc("/chat", websocket.ChatHandler)

	// 消息历史端点
	r.HandleFunc("/message/list", func(w http.ResponseWriter, r *http.Request) {
		db := database.InitDB()
		defer func(db *database.DB) {
			err := db.Close()
			if err != nil {

			}
		}(db)

		customerID := r.URL.Query().Get("customer_id")
		var messages []models.Message
		customerIDUint, err := strconv.ParseUint(customerID, 10, 32)
		if err != nil {
			http.Error(w, `{"error": "Invalid customer_id"}`, http.StatusBadRequest)
			return
		}
		db.Debug().Table("messages").Where("customer_id = ?", uint(customerIDUint)).Find(&messages)
		json.NewEncoder(w).Encode(messages)
	}).Methods("GET")

	return r
}

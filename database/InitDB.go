package database

import (
	"github.com/jinzhu/gorm"
	"github.com/larking7days/websocket_chatbot/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*gorm.DB
}

func InitDB() *DB {
	db, err := gorm.Open("sqlite3", "chat.db") // 修正数据库文件名扩展
	if err != nil {
		panic("数据库连接失败")
	}
	db.AutoMigrate(&models.Customer{}, &models.Message{}, &models.Feedback{})
	return &DB{db}
}

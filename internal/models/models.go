package models

import "time"

type Customer struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
}

type Message struct {
	ID         uint `gorm:"primaryKey"`
	CustomerID uint
	Content    string
	SenderType string // "customer" 或 "bot"
	CreatedAt  time.Time
	Sentiment  string // 新增情感分析字段
	Rating     string `json:"rating" gorm:"size:1"` // 新增评分字段
}

type Feedback struct {
	ID          uint `gorm:"primaryKey"`
	CustomerID  uint
	Rating      int // 1-5
	Comment     string
	RespondedTo bool // 标记是否已回复
	CreatedAt   time.Time
	Sentiment   string `gorm:"index"` // 添加索引提升查询性能
}

// 新增AI处理结构体
type AIConfig struct {
	OpenAIToken  string `yaml:"-"` // 忽略直接配置
	AnthropicKey string `yaml:"-"`
	APIToken     string `yaml:"api_token"` // 通过环境变量注入
}

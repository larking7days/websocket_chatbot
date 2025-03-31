package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/larking7days/websocket_chatbot/database"
	"github.com/larking7days/websocket_chatbot/internal/ai"
	. "github.com/larking7days/websocket_chatbot/internal/models"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// 在ChatHandler顶部添加
const (
	SentimentPositive = "positive"
	SentimentNeutral  = "neutral"
	SentimentNegative = "negative"
)

// 在ChatHandler函数开头添加
var (
	aiAnalyzer *ai.Analyzer // 新增AI分析器
	chatConfig map[string]interface{}
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 添加跨域支持（生产环境需限制）
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	// 加载聊天规则
	configFile, err := os.ReadFile("config/chat_rules.yaml")
	if err != nil {
		panic("加载配置文件失败: " + err.Error())
	}
	if err := yaml.Unmarshal(configFile, &chatConfig); err != nil {
		panic("解析配置文件失败: " + err.Error())
	}

	aiAnalyzer = ai.NewAnalyzer()
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		// 移除 http.Error 调用，避免重复写响应头
		println("WebSocket 升级失败详情:",
			"\n错误信息:", err.Error(),
			"\n请求头:", r.Header)
		return
	}
	defer conn.Close()
	db := database.InitDB()
	defer db.Close()
	// 添加连接状态跟踪
	clientAddr := conn.RemoteAddr().String()
	log.Printf("clinet conn: %s", clientAddr)
	defer log.Printf("client disconnect: %s", clientAddr)

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			// 添加连接关闭处理
			if websocket.IsUnexpectedCloseError(err) {
				println("conn disconnect:", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("Failed to parse message [客户端:%s]: %v", clientAddr, err) // 改用标准日志
			continue
		}
		if msg.SenderType == "" {
			msg.SenderType = "customer"
		}
		// 处理消息并保存到数据库
		if err := processMessage(db, conn, msg); err != nil {
			logError("消息处理失败:", err)
			continue
		}

	}

}

// 新增消息处理逻辑
func processMessage(db *database.DB, conn *websocket.Conn, msg Message) error {
	// 执行情感分析
	if aiAnalyzer != nil {
		if sentiment, err := aiAnalyzer.ClassifySentiment(msg.Content); err == nil {
			msg.Sentiment = sentiment
		}
	}

	// 保存原始消息
	if err := db.Create(&msg).Error; err != nil {
		return err
	}
	// 新增评分处理逻辑
	if isRatingResponse(msg) {
		return handleRating(db, conn, msg)
	}
	// 处理工作流逻辑
	if response := handleWorkflow(msg); response != "" {
		return sendBotResponse(db, conn, msg.CustomerID, response)
	}

	// 处理反馈逻辑
	if shouldTrigger(msg) {
		return triggerFeedbackFlow(db, conn, msg.CustomerID)
	}

	return nil
}

// 新增评分处理函数
func isRatingResponse(msg Message) bool {
	// 仅处理用户发送的数字评分
	return msg.SenderType == "customer" &&
		len(strings.TrimSpace(msg.Content)) == 1 &&
		strings.ContainsAny(msg.Content, "12345")
}
func handleRating(db *database.DB, conn *websocket.Conn, msg Message) error {
	// 记录评分并发送确认
	rating := strings.TrimSpace(msg.Content)
	responseText := fmt.Sprintf("thanks for rating：%s！", rating)

	// 保存评分更新
	db.Model(&msg).Update("rating", rating)

	return sendBotResponse(db, conn, msg.CustomerID, responseText)
}

// 新增工作流处理函数
func handleWorkflow(msg Message) string {
	rules, ok := chatConfig["workflows"].(map[string]interface{})
	if !ok {
		log.Printf("无效的workflows配置格式")
		return ""
	}

	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})
		keywords := convertToStringSlice(ruleMap["trigger_keywords"].([]interface{}))

		// 关键词匹配
		for _, kw := range keywords {
			pattern := `\b` + regexp.QuoteMeta(kw) + `\b`
			if matched, _ := regexp.MatchString(pattern, strings.ToLower(msg.Content)); matched {
				return ruleMap["response"].(string)
			}
		}

		// 情感匹配
		if sentimentRules, ok := ruleMap["sentiment_triggers"].([]interface{}); ok {
			for _, s := range sentimentRules {
				if strings.ToLower(msg.Sentiment) == strings.ToLower(s.(string)) {
					return ruleMap["response"].(string)
				}
			}
		}
	}
	// 新增AI增强响应（当没有匹配规则时）
	if aiAnalyzer != nil {
		enhancedResponse, err := aiAnalyzer.EnhanceResponse(msg.Content, msg.Sentiment)
		if err == nil {
			return enhancedResponse
		}
	}

	return ""
}

// 新增通用工具函数
func convertToStringSlice(input []interface{}) []string {
	result := make([]string, len(input))
	for i, v := range input {
		result[i] = v.(string)
	}
	return result
}

func sendBotResponse(db *database.DB, conn *websocket.Conn, customerID uint, content string) error {
	response := Message{
		CustomerID: customerID,
		Content:    content,
		SenderType: "bot",
	}

	if err := db.Create(&response).Error; err != nil {
		return err
	}

	return conn.WriteJSON(response)
}

func triggerFeedbackFlow(db *database.DB, conn *websocket.Conn, customerID uint) error {
	feedbackMsg := Message{
		CustomerID: customerID,
		Content:    "Please rate your experience from 1-5:",
		SenderType: "bot",
	}

	if err := db.Create(&feedbackMsg).Error; err != nil {
		return err
	}

	return conn.WriteJSON(feedbackMsg)
}

// 修改后的触发条件判断
func shouldTrigger(msg Message) bool {
	// 原有关键词触发
	if shouldTriggerFeedback(msg.Content) {
		return true
	}

	// 新增负面情感触发
	if msg.Sentiment == "negative" {
		return true
	}

	return false
}

// 原有工具函数保持不变
func shouldTriggerFeedback(text string) bool {
	keywords := []string{"feedback", "review"}
	lowerText := strings.ToLower(text)

	// 使用正则表达式进行全词匹配
	for _, kw := range keywords {
		pattern := `\b` + regexp.QuoteMeta(kw) + `\b`
		if matched, _ := regexp.MatchString(pattern, lowerText); matched {
			return true
		}
	}
	return false
}

// 新增错误处理函数
func handleConnectionError(err error) {
	if websocket.IsUnexpectedCloseError(err) {
		println("客户端异常断开:", err)
	}
}

func logError(prefix string, err error) {
	println(prefix, err.Error())
}

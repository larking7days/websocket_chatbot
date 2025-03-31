# AI-Powered Conversational Service Platform


Real-time intelligent conversation system with sentiment analysis and automated workflows

## ✨ Core Features

- **Sentiment Recognition**: Real-time emotion detection (Positive/Neutral/Negative)
- **Smart Routing**: Context-aware conversation flows
- **Multi-AI Support**: Alibaba Qwen & OpenAI integration
- **Conversation Persistence**: Full chat history storage
- **Bi-directional Communication**: WebSocket real-time messaging

## 🛠 Technical Architecture

```plaintext
+------------------+
|  Web Client       |
+--------+---------+
         |
+--------+---------+
| WebSocket Gateway | (Gorilla)
+--------+---------+
         |
+--------+---------+
|  AI Processing   | (Qwen/OpenAI)
+--------+---------+
         |
+--------+---------+
| Data Persistence | (GORM + SQLite)
+------------------+
```
## 🚀 Getting Started

```

## Prerequisites
Go 1.20+
SQLite3
Aliyun DashScope API Key

1. **Clone Repository**: `git clone URL_ADDRESS1. **Clone Repository**: `git clone https://github.com/lang7days/websocket_ai.git`
2. **Install Dependencies**: `go mod download`
3. **export DASHSCOPE_API_KEY=your_api_key_here
4  **Edit config/chat_rules.yaml to modify conversation workflows:   
5. **Run the server**: `go run cmd/server/main.go`
6. **Run the client**: `go run cmd/cli/main.go`
```
### ⚙️ Workflow Automation
```yaml
# Sample Workflow Rule
#config/chat_rules.yaml
workflows:
  feedback:
    trigger_keywords: ["feedback", "review"]
    response: "Please rate your experience from 1-5. 1 is the worst, 5 is the best."
    next_step: "collect_rating"
```

## 💡 Usage Examples

### Basic Interaction
```
// Client sends:
You:
{"Content": "Hello", "CustomerId": 123}

// Server responds: 
 Bot:
{"ID":2,"CustomerID":123,"Content":"Hey there! How's it going?","SenderType":"bot","CreatedAt":"2025-03-31T12:57:23.455219+08:00","Sentiment":"","rating":""}

// Client sends:
You:
{"Content": "preview", "CustomerId": 123}

// Server responds: 
Bot:
 {"ID":4,"CustomerID":123,"Content":"当然！可以这样回复用户：\n\n\"嘿，preview听起来很棒！你是在分享一些有趣的东西吗？我很期待能看到或了解到更多呢！如果有任何需要帮忙的地方，尽管告诉我哦！\"","SenderType":"bot","CreatedAt":"2025-03-31T12:58:09.629125+08:00","Sentiment":"","rating":""}

You: 
{"Content": "review", "CustomerId": 123}
Bot:
 {"ID":6,"CustomerID":123,"Content":"Please rate your experience from 1-5. 1 is the worst, 5 is the best.","SenderType":"bot","CreatedAt":"2025-03-31T12:59:10.6321+08:00","Sentiment":"","rating":""}
You:
{"Content": "5", "CustomerId": 123} 
Bot:
 {"ID":8,"CustomerID":123,"Content":"thanks for rating：5！","SenderType":"bot","CreatedAt":"2025-03-31T12:59:43.049993+08:00","Sentiment":"","rating":""}

```
### Retrieves chat history for a customer.

GET /message/list 
````
http://localhost:8080/message/list?customer_id=123

```json
[{"ID":1,"CustomerID":123,"Content":"Hello","SenderType":"customer","CreatedAt":"2025-03-31T12:57:19.374135+08:00","Sentiment":"neutral","rating":""},
{"ID":2,"CustomerID":123,"Content":"Hey there! How's it going?","SenderType":"bot","CreatedAt":"2025-03-31T12:57:23.455219+08:00","Sentiment":"","rating":""},
{"ID":3,"CustomerID":123,"Content":"preview","SenderType":"customer","CreatedAt":"2025-03-31T12:58:03.78543+08:00","Sentiment":"positive","rating":""},
{"ID":4,"CustomerID":123,"Content":"当然！可以这样回复用户：\n\n\"嘿，preview听起来很棒！你是在分享一些有趣的东西吗？我很期待能看到或了解到更多呢！如果有任何需要帮忙的地方，尽管告诉我哦！\"","SenderType":"bot","CreatedAt":"2025-03-31T12:58:09.629125+08:00","Sentiment":"","rating":""},
{"ID":5,"CustomerID":123,"Content":"review","SenderType":"customer","CreatedAt":"2025-03-31T12:59:08.53468+08:00","Sentiment":"positive","rating":""},{"ID":6,"CustomerID":123,"Content":"Please rate your experience from 1-5. 1 is the worst, 5 is the best.","SenderType":"bot","CreatedAt":"2025-03-31T12:59:10.6321+08:00","Sentiment":"","rating":""},
{"ID":7,"CustomerID":123,"Content":"5","SenderType":"customer","CreatedAt":"2025-03-31T12:59:42.404157+08:00","Sentiment":"neutral","rating":"5"},
{"ID":8,"CustomerID":123,"Content":"thanks for rating：5！","SenderType":"bot","CreatedAt":"2025-03-31T12:59:43.049993+08:00","Sentiment":"","rating":""}]


```
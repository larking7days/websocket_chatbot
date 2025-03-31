package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
)

func main() {
	// 添加错误处理三件套
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/chat", nil)
	if err != nil {
		fmt.Printf("连接失败: %v\n", err) // 显示具体错误信息
		os.Exit(1)
	}
	defer conn.Close() // 添加关闭连接 defer

	go func() {
		for {
			// 添加读消息错误处理
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("\n连接中断:", err)
				os.Exit(1)
			}
			fmt.Println("\nBot:", string(message))
			fmt.Print("You: ")
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You: ")
		text, _ := reader.ReadString('\n')
		// 添加写消息错误处理
		if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			fmt.Println("发送失败:", err)
			os.Exit(1)
		}
	}
}

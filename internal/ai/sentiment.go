package ai

import (
	"context"
	"errors"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	_ "github.com/openai/openai-go/option"
	"os"
	"strings"
)

type Analyzer struct {
	client *openai.Client
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		client: openai.NewClient(
			option.WithAPIKey(os.Getenv("DASHSCOPE_API_KEY")),
			option.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1/"),
		),
	}
}

func (a *Analyzer) ClassifySentiment(text string) (string, error) {
	resp, err := a.client.Chat.Completions.New(
		context.TODO(), openai.ChatCompletionNewParams{
			Messages: openai.F(
				[]openai.ChatCompletionMessageParamUnion{
					openai.UserMessage("Classify feedback as positive, neutral, or negative: " + text),
				},
			),
			Model: openai.F("qwen-plus"),
		},
	)

	if err != nil {
		panic(err.Error())
	}

	// 简化的情感解析逻辑
	if err == nil && len(resp.Choices) > 0 {
		return parseSentiment(resp.Choices[0].Message.Content)
	}
	return "neutral", err
}

func parseSentiment(response string) (string, error) {
	// 实际实现需要处理AI返回的自然语言响应
	// 这里使用简化版本演示
	// Normalize response for easier processing
	lowerResponse := strings.ToLower(response)

	switch {
	case strings.Contains(lowerResponse, "positive**"),
		strings.Contains(lowerResponse, "good**"),
		strings.Contains(lowerResponse, "excellent**"):
		return "positive", nil
	case strings.Contains(lowerResponse, "negative**"),
		strings.Contains(lowerResponse, "poor**"),
		strings.Contains(lowerResponse, "bad**"):
		return "negative", nil
	default:
		// Fallback to neutral for ambiguous cases
		return "neutral", nil
	}
}
func (a *Analyzer) EnhanceResponse(query string, sentiment string) (string, error) {
	// 示例提示词工程
	prompt := fmt.Sprintf(`基于用户情绪(%s)生成友好回复：
    用户最后消息：%s
    生成自然的口语化响应：`, sentiment, query)

	// 调用AI模型生成
	return a.generateResponse(prompt)
}
func (a *Analyzer) generateResponse(prompt string) (string, error) {
	resp, err := a.client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			}),
			Model:       openai.F("qwen-plus"),
			Temperature: openai.F(0.7), // 新增创造性控制参数
		},
	)

	if err != nil {
		return "", fmt.Errorf("AI响应生成失败: %w", err)
	}

	if len(resp.Choices) > 0 {
		return strings.TrimSpace(resp.Choices[0].Message.Content), nil
	}
	return "", errors.New("未生成有效响应")
}

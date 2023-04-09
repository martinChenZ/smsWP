package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"strings"
	"time"
)

type ChoiceMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GptTurbo struct {
	model    string
	messages []ChoiceMsg
	user     string
}

var Appkey string

func CallGpt3(msg []ChoiceMsg) string {
	fmt.Println(Appkey)

	gpturl := "https://api.openai.com/v1/chat/completions"
	message := make([]ChoiceMsg, 1)
	turbo := &GptTurbo{
		"gpt-3.5-turbo",
		message,
		"1",
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	data, _ := json.Marshal(turbo)
	req, err := http.NewRequest("post", gpturl, strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("Error send to gpt ", err)
		return "请求失败"
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", Appkey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error send to gpt ", err)
		return "请求失败"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading  resp ", err)
		return "请求失败"
	}

	return string(body)
}

func CallGpt(msg []openai.ChatCompletionMessage) string {
	client := openai.NewClient(Appkey)
	if msg[0].Role != openai.ChatMessageRoleSystem {
		msg = append([]openai.ChatCompletionMessage{
			{Role: "system",
				Content: "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible."}}, msg...)
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: msg,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "请求失败"
	}

	fmt.Println(resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content
}

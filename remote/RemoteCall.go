package remote

import (
	"encoding/json"
	"fmt"
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

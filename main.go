package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	fmt.Println("chat：你好！我是zq的ai助手！有什么可以帮助您的吗？")
	var parentMessageId string
	for {
		var input string
		fmt.Print("> ")
		fmt.Scanln(&input)
		url := "https://www.muzjia.com/api/chat-process"
		// 创建请求体
		options := map[string]string{
			"parentMessageId": parentMessageId,
		}
		body := map[string]interface{}{
			"prompt":  input,
			"options": options,
		}
		requestBody, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		messages := httpPost(url, requestBody)
		parentMessageId = messages[len(messages)-1].ParentMessageId
		fmt.Println("chat：", messages[len(messages)-1].Text)
	}
}

type Message struct {
	ParentMessageId string `json:"parentMessageId"`
	Text            string `json:"text"`
}

// 发送 POST 请求
func httpPost(url string, requestBody []byte) []Message {
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	bodyJson := "[" + string(bodyBytes)
	bodyJson = string(bodyJson) + "]"
	bodyJson = strings.Replace(bodyJson, "}\n{", "},{", -1)
	bodyJson = strings.Replace(bodyJson, "}{", "},{", -1)
	// 处理响应数据
	var messages []Message
	err = json.Unmarshal([]byte(bodyJson), &messages)
	if err != nil {
		panic(err)
	}

	return messages
}

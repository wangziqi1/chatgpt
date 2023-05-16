package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("chat：你好！我是zq的ai助手！有什么可以帮助您的吗？")
	fmt.Print("\n")
	var parentMessageId string
	for {
		var input string
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("你：")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		input = strings.TrimSpace(input)
		fmt.Print("\n")
		sig := make(chan bool)
		go spinner(time.Duration(time.Millisecond*40), sig)
		// fmt.Scanln(&input)
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
		sig <- true
		fmt.Println("\rchat：", messages[len(messages)-1].Text)
		fmt.Print("\n")
	}
}

func spinner(delay time.Duration, done chan bool) {
	signs := [6]string{".     ", "..    ", "...   ", "....  ", "..... ", "......"}
	for {
		for _, i := range signs {
			select {
			case <-done: //通道接收到信号退出
				return
			case <-time.After(delay): //在一定时间内没有收到信号
				fmt.Printf("\rchat：思考中%s", i)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

type Message struct {
	ParentMessageId string `json:"parentMessageId"`
	Text            string `json:"text"`
}

// 发送 POST 请求
func httpPost(url string, requestBody []byte) []Message {
	// Create a new HTTP client with InsecureSkipVerify set to true
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Post(url, "application/json", bytes.NewBuffer(requestBody))
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

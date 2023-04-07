package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"bytes"
)

type OpenAIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Error *OpenAIError `json:"error"`
}

func main() {
	// コマンドライン引数を確認し、プロンプトを取得します。
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [your-text]")
		os.Exit(1)
	}
	prompt := os.Args[1]

	// APIを叩きます。
	body, err := getAPIResponse(prompt)
	if err != nil {
		log.Fatalf("Error calling the API: %v", err)
	}

	// APIレスポンスを表示します。
	displayAPIResponse(body)
}

func processArguments() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [your-text]")
		os.Exit(1)
	}
}

func getAPIResponse(prompt string) ([]byte, error) {
	apiKey, err := loadAPIKey()
	if err != nil {
		return nil, err
	}

	params := buildAPIParams(prompt)
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	endpoint := "https://api.openai.com/v1/chat/completions"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(paramsJSON))
	if err != nil {
		return nil, err
	}

	setRequestHeaders(req, apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func displayAPIResponse(body []byte) {
	fmt.Printf("Full API response: %s\n", string(body))

	var apiResponse OpenAIResponse
	err := json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling the API response: %v", err)
	}

	if &apiResponse != nil {
		if apiResponse.Error != nil {
			fmt.Printf("Error from the OpenAI API: %s - %s\n", apiResponse.Error.Code, apiResponse.Error.Message)
		} else if len(apiResponse.Choices) > 0 {
			fmt.Printf("Response from the OpenAI API: %s\n", apiResponse.Choices[0].Message.Content)
		} else {
			fmt.Println("No choices returned from the OpenAI API")
		}
	}
}

func buildAPIParams(prompt string) map[string]interface{} {
	return map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 1.0,
	}
}

func setRequestHeaders(req *http.Request, apiKey string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
}

func loadAPIKey() (string, error) {
	apiKeyBytes, err := ioutil.ReadFile("api_key.txt")
	if err != nil {
		return "", err
	}
	return string(apiKeyBytes), nil
}

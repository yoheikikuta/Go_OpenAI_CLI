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
        Message struct { // 追加
            Role    string `json:"role"`
            Content string `json:"content"` // TextからContentに変更
        } `json:"message"`
        Index        int         `json:"index"`
        Logprobs     interface{} `json:"logprobs"`
        FinishReason string      `json:"finish_reason"`
    } `json:"choices"`
    Error *OpenAIError `json:"error"`
}

func loadAPIKey() (string, error) {
	apiKeyBytes, err := ioutil.ReadFile("api_key.txt")
	if err != nil {
		return "", err
	}
	return string(apiKeyBytes), nil
}

func callOpenAIAPI(prompt string) ([]byte, error) {
	// APIエンドポイントと認証情報を設定します。
  apiKey, err := loadAPIKey()
	endpoint := "https://api.openai.com/v1/chat/completions"

    // APIリクエストを作成します。
    params := map[string]interface{}{
        "model":      "gpt-3.5-turbo",
        "messages": []map[string]string{
            {
                "role": "user",
                "content": prompt,
            },
        },
        "temperature": 1.0,
    }

    paramsJSON, err := json.Marshal(params)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(paramsJSON))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

    // APIリクエストを実行します。
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // APIレスポンスを読み込みます。
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
}

func processArguments() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [your-text]")
		os.Exit(1)
	}
}

func main() {
	processArguments()

	// APIを叩く
	prompt := "What is the capital of France?"
	body, err := callOpenAIAPI(prompt)
	if err != nil {
		log.Fatalf("Error calling the API: %v", err)
	}

	// APIレスポンス全体を出力
	fmt.Printf("Full API response: %s\n", string(body))

	// APIレスポンスをデコードします。
	var apiResponse OpenAIResponse
	err = json.Unmarshal(body, &apiResponse)
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


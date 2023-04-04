package main

import (
	"fmt"
	"log"
	"os"
)

type MockAPIResponse struct {
	Message string `json:"message"`
}

func callMockAPI() (*MockAPIResponse, error) {
	// ここで実際のAPIエンドポイントと認証情報を使用します。
	response := &MockAPIResponse{
		Message: "Hello from the Mock OpenAI API!",
	}
	return response, nil
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
	response, err := callMockAPI()
	if err != nil {
		log.Fatalf("Error calling the API: %v", err)
	}

	// 結果を表示する
	fmt.Printf("Response from the Mock OpenAI API: %s\n", response.Message)
}

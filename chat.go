package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
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
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Error *OpenAIError `json:"error"`
}

var (
	// ANSI escape codes for text colors
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	messages := make([]map[string]string, 0)

	var wg sync.WaitGroup

	for {
		fmt.Print("Enter your message: ")
		prompt, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}

		prompt = strings.TrimSpace(prompt)
		if prompt == "exit" {
			break
		}

		userMessage := map[string]string{
			"role":    "user",
			"content": prompt,
		}
		messages = append(messages, userMessage)

		params, err := buildAPIParams(prompt, messages)
		if err != nil {
			log.Fatalf("Error building API params: %v", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := getAPIResponse(params)
			if err != nil {
				log.Fatalf("Error calling the API: %v", err)
			}
		}()

		wg.Wait()
	}
}

func getAPIResponse(params map[string]interface{}) error {
	apiKey, err := loadAPIKey()
	if err != nil {
		return err
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return err
	}

	endpoint := "https://api.openai.com/v1/chat/completions"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(paramsJSON))
	if err != nil {
		return err
	}

	setRequestHeaders(req, apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// 空行を読み飛ばす
		if strings.TrimSpace(line) == "" {
			continue
		}

		// fmt.Println("Raw response:", line)  // for debugging.

		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		}

		if strings.TrimSpace(line) == "[DONE]" {
			// Stream has ended.
			fmt.Println("\n")
			break
		}

		var apiResponse OpenAIResponse
		err = json.Unmarshal([]byte(line), &apiResponse)
		if err != nil {
			log.Fatal(err)
		}

		if apiResponse.Error != nil {
			fmt.Printf("Error from the OpenAI API: %s - %s\n", apiResponse.Error.Code, apiResponse.Error.Message)
		} else if len(apiResponse.Choices) > 0 {
			deltaContent := apiResponse.Choices[0].Delta.Content
			if deltaContent != "" {
				fmt.Printf("%s%s%s", colorCyan, deltaContent, colorReset)
			}
		} else {
			fmt.Println("No choices returned from the OpenAI API")
		}
	}

	return nil
}

func buildAPIParams(prompt string, messages []map[string]string) (map[string]interface{}, error) {
	if len(messages) == 0 {
		return nil, errors.New("messages should not be empty")
	}

	messages = append(messages, map[string]string{
		"role":    "user",
		"content": prompt,
	})

	maxTokens := 500
	temperature := 0.8
	topP := 0.9

	return map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"messages":    messages,
		"max_tokens":  maxTokens,
		"temperature": temperature,
		"top_p":       topP,
		"stream":      true,
	}, nil
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

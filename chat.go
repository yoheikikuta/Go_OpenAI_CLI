package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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
    reader := bufio.NewReader(os.Stdin)
    messages := make([]map[string]string, 0)

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

        params := buildAPIParams(prompt, messages)
        body, err := getAPIResponse(params)
        if err != nil {
            log.Fatalf("Error calling the API: %v", err)
        }

        responseMessage := displayAPIResponse(body)
        messages = append(messages, responseMessage)
    }
}

func getAPIResponse(params map[string]interface{}) ([]byte, error) {
	apiKey, err := loadAPIKey()
	if err != nil {
		return nil, err
	}

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

func displayAPIResponse(body []byte) map[string]string {
	var apiResponse OpenAIResponse
	err := json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling the API response: %v", err)
	}

	if &apiResponse != nil {
		if apiResponse.Error != nil {
			fmt.Printf("Error from the OpenAI API: %s - %s\n", apiResponse.Error.Code, apiResponse.Error.Message)
		} else if len(apiResponse.Choices) > 0 {
			response := apiResponse.Choices[0].Message.Content
			fmt.Printf("Response from the OpenAI API: %s\n", response)
			return map[string]string{
				"role":    "assistant",
				"content": response,
			}
		} else {
			fmt.Println("No choices returned from the OpenAI API")
		}
	}
	return nil
}

func buildAPIParams(prompt string, messages []map[string]string) map[string]interface{} {
    messages = append(messages, map[string]string{
        "role":    "user",
        "content": prompt,
    })

    maxTokens := 50
    temperature := 0.8
    topP := 0.9

    return map[string]interface{}{
        "model":       "gpt-3.5-turbo",
        "messages":    messages,
        "max_tokens":  maxTokens,
        "temperature": temperature,
        "top_p":       topP,
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
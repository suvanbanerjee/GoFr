package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func predictSentiment(text string) (string, error) {
	apiURL := "https://api-inference.huggingface.co/models/tabularisai/robust-sentiment-analysis"
	// Replace with your Hugging Face API key
	apiKey := "your_hugging_face_api_key"

	// Prepare request data
	payload := map[string]interface{}{
		"inputs": text,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create a request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the response JSON
	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Extract sentiment class from the response
	if len(result) > 0 {
		if score, exists := result[0]["label"]; exists {
			return score.(string), nil
		}
	}

	return "", fmt.Errorf("sentiment not found in response")
}

func main() {
	texts := []string{
		"I absolutely loved this movie! The acting was superb and the plot was engaging.",
		"The service at this restaurant was terrible. I'll never go back.",
		"The product works as expected. Nothing special, but it gets the job done.",
		"I'm somewhat disappointed with my purchase. It's not as good as I hoped.",
		"This book changed my life! I couldn't put it down and learned so much.",
	}

	for _, text := range texts {
		sentiment, err := predictSentiment(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Text: %s\nSentiment: %s\n\n", text, sentiment)
	}
}

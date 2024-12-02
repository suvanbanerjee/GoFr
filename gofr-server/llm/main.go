package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const ollamaAPI = "http://localhost:11434/api/generate" // Example API endpoint

func Llm() string {
	// Read the files and concatenate them
	// concatenatedFilePath := "/Users/yuktha/Documents/Go idk/output.txt" // Replace with the actual path

	// Get the context from the concatenated file
	fileContent, err := ioutil.ReadFile("/Users/yuktha/Downloads/GoFr Server 2/output.txt")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// fmt.Print("Context:", string(fileContent))
	// Create a sample request payload to send to TinyLlama
	requestPayload := map[string]interface{}{
		"model":  "tinyllama",
		"prompt": "Here is the content of the text file:" + string(fileContent) + "Please only respond with information derived from the above content. \nmake a Twitter post about Custom Spans In Tracing",
		"stream": false,
	}
	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Fatalf("Error marshalling request: %v", err)
	}

	// Send the request to the Ollama API
	resp, err := http.Post(ollamaAPI, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	// Print the response
	// fmt.Println("Response from TinyLlama:", response)
	if respText, exists := response["response"].(string); exists {
		fmt.Println("Response from TinyLlama:", respText)
		return respText
	} else {
		fmt.Println("Response field not found in the API response")
	}
	return "Error"
}

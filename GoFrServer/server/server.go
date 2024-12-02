package main

import (
	"GoFr/GoFrServer/llm"
	"GoFr/GoFrServer/sendmail"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"gofr.dev/pkg/gofr"
)

type LLMRequest struct {
	Question string `json:"question"`
}

type LLMResponse struct {
	Status   string `json:"status"`
	Response string `json:"response"`
}

func callLLMAPI(context string) (string, error) {
	apiURL := "http://localhost:8000/generate_content/"

	requestBody := LLMRequest{
		Question: context,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var llmResponse LLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&llmResponse); err != nil {
		return "", err
	}

	return llmResponse.Response, nil
}

func main() {
	err := godotenv.Load("GoFrServer/configs/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	app := gofr.New()

	// Chain with Sentiment Analysis API
	app.GET("/get_twit_trend", func(ctx *gofr.Context) (interface{}, error) {
		// Define the Python script to execute
		cmd := exec.Command("python", "twittrend.py")

		// Capture the output
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		// Run the Python script
		err := cmd.Run()
		if err != nil {
			// Return any error during script execution
			return nil, errors.New("Failed to run Python script: " + stderr.String())
		}

		// Get the output text from the script
		result := strings.TrimSpace(out.String())

		// Check if the result is empty
		if result == "" {
			return nil, errors.New("No tweets found or script output is empty")
		}

		// Return the concatenated tweets
		return map[string]string{
			"tweets": result,
		}, nil
	})

	app.POST("/create/x", func(ctx *gofr.Context) (interface{}, error) {
		// Get the text from request body
		var body struct {
			Text string `json:"text"`
		}
		if err := ctx.Bind(&body); err != nil {
			return nil, errors.New("invalid request body")
		}

		if body.Text == "" {
			return nil, errors.New("text is required")
		}

		Restes := llm.Llm()

		ctx.Logger.Info("Received text: ", Restes)

		cmd := exec.Command("python", "twit.py", body.Text)
		err := cmd.Run()
		if err != nil {
			return nil, errors.New("failed to run the Python script")
		}

		return map[string]string{
			"message": "Text received successfully and script executed",
		}, nil
	})

	app.POST("/create/email", func(ctx *gofr.Context) (interface{}, error) {
		var body struct {
			Context string `json:"context"`
		}
		if err := ctx.Bind(&body); err != nil {
			return nil, errors.New("invalid request body")
		}

		if body.Context == "" {
			return nil, errors.New("context is required")
		}

		// Get the LLM-generated content
		generatedContent, err := callLLMAPI(body.Context)
		if err != nil {
			return nil, errors.New("failed to generate content: " + err.Error())
		}

		ctx.Logger.Info("Generated content: ", generatedContent)

		// Send the generated content via email
		sendmail.Send_mail(generatedContent)

		return map[string]string{
			"message": "Email generated and sent successfully",
			"content": generatedContent,
		}, nil
	})

	app.Run()
}

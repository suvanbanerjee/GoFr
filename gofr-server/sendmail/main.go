package sendmail

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Get_llm_response makes an API request to generate email content.
func Get_llm_response(context string) (string, string) {
	template := `You are an email generator tasked with creating engaging, concise, and professional social media posts.

	Use only the information provided in the context.
	Do not add any extra details, assumptions, or speculative content.
	Maintain a tone suitable for [platform: e.g., LinkedIn, Twitter, Instagram, etc.].
	The post should be clear, concise, and adhere to any specified character limits or formatting guidelines.
	If the content is technical or professional, ensure the language is precise and jargon-free (if possible). For creative or casual posts, keep the tone friendly and approachable.
	`

	// Construct the full prompt
	prompt := context + template
	url := "http://localhost:8000/generate_post" // Replace with correct URL

	// Creating JSON payload
	jsonStr := []byte(`{"Context":"` + prompt + `"}`)

	// Making the POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatalf("Failed to create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %s", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %s", err)
	}

	// Check if the response status code indicates success
	if resp.StatusCode != 200 {
		log.Printf("Error from API: %d - %s", resp.StatusCode, string(body))
		return "_ERROR", "_ERROR"
	}

	// Log the full response for debugging
	log.Printf("Response Body: %s", string(body))

	// Assuming the response body is a JSON with 'subject' and 'body'
	// This part can be modified depending on the exact format of the response
	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Failed to parse JSON response: %s", err)
	}

	return result["body"], result["subject"]
}

// Send_mail sends the generated email to a list of recipients from a CSV file.
func Send_mail(context string) {
	// Open the CSV file containing emails
	file, err := os.Open("emails.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Read all email records
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	// Get the subject and body from LLM response
	text, title := Get_llm_response(context)
	if text == "_ERROR" || title == "_ERROR" {
		log.Fatalf("Failed to get response")
	} else {
		for _, record := range records {
			toEmail := record[0] // Email address from the CSV file
			from := mail.NewEmail("Suvan Banerje", "suvan@burdenoff.com")
			subject := title
			to := mail.NewEmail("User", toEmail)
			plainTextContent := text
			htmlContent := "<strong>" + text + "</strong>"

			// Create a SendGrid message
			message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

			// Send email using SendGrid client
			client := sendgrid.NewSendClient("apikey")
			response, err := client.Send(message)
			if err != nil {
				log.Println(err)
			} else {
				// Log the SendGrid response
				log.Printf("SendGrid Response: StatusCode: %d, Body: %s, Headers: %s",
					response.StatusCode, response.Body, response.Headers)
			}
		}
	}
}

// GetAnalytics fetches the analytics data from SendGrid.
func GetAnalytics() {
	apiKey := "apikey" // Replace with correct API key
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(apiKey, "/v3/stats", host)
	request.Method = "GET"
	queryParams := make(map[string]string)
	queryParams["start_date"] = "2024-11-22"
	request.QueryParams = queryParams

	// Make the API request
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		// Log the analytics response
		log.Printf("SendGrid Analytics: StatusCode: %d, Body: %s, Headers: %s",
			response.StatusCode, response.Body, response.Headers)
	}
}

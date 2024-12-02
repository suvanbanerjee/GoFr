package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GitHub repository details
const (
	owner      = "gofr-dev"
	repo       = "gofr"
	basePath   = "docs"     // The folder in the repo to start scraping
	branch     = "development"
	outputDir  = "./gogo"   // Directory to save markdown files
	apiBaseURL = "https://api.github.com/repos/" + owner + "/" + repo + "/contents"
)

// File represents a GitHub API response item
type File struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

// getFilesRecursively fetches file metadata recursively from the GitHub API
func getFilesRecursively(path string) ([]string, error) {
	url := fmt.Sprintf("%s/%s?ref=%s", apiBaseURL, path, branch)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching %s: %s", url, resp.Status)
	}

	var items []File
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var files []string
	for _, item := range items {
		if item.Type == "file" && strings.HasSuffix(item.Name, ".md") {
			files = append(files, item.DownloadURL)
		} else if item.Type == "dir" {
			subFiles, err := getFilesRecursively(item.Path)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		}
	}

	return files, nil
}

// downloadFiles downloads markdown files to the output directory with unique filenames
func downloadFiles(links []string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory %s: %v", outputDir, err)
		}
	}

	for _, link := range links {
		resp, err := http.Get(link)
		if err != nil {
			fmt.Printf("Error downloading %s: %v\n", link, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error downloading %s: %s\n", link, resp.Status)
			continue
		}

		originalFileName := filepath.Base(link)
		timestamp := time.Now().Unix()
		uniqueFileName := fmt.Sprintf("%s_%d%s",
			strings.TrimSuffix(originalFileName, filepath.Ext(originalFileName)),
			timestamp,
			filepath.Ext(originalFileName),
		)

		filePath := filepath.Join(outputDir, uniqueFileName)
		fmt.Printf("Downloading %s to %s\n", link, filePath)

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", filePath, err)
			continue
		}
		defer file.Close()

		if _, err := io.Copy(file, resp.Body); err != nil {
			fmt.Printf("Error writing to file %s: %v\n", filePath, err)
		}
	}

	return nil
}

func main() {
	fmt.Println("Starting scrape...")
	markdownLinks, err := getFilesRecursively(basePath)
	if err != nil {
		fmt.Printf("Error fetching files: %v\n", err)
		return
	}

	fmt.Printf("Found %d markdown files.\n", len(markdownLinks))
	if err := downloadFiles(markdownLinks); err != nil {
		fmt.Printf("Error downloading files: %v\n", err)
		return
	}

	fmt.Println("Download complete!")
}

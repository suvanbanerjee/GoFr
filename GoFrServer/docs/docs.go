package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	// Specify the directory containing .md files
	dir := "/Users/yuktha/Downloads/gofr-development/docs" // Set your directory here

	// Create the output file
	outputFile, err := os.Create("combined.md")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Read all files in the directory
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the file is a .md file
		if filepath.Ext(path) == ".md" {
			// Read the content of the file
			content, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return err
			}

			// Write the content to the output file
			_, err = outputFile.Write(content)
			if err != nil {
				fmt.Println("Error writing to output file:", err)
				return err
			}

			// Add a newline between files
			_, err = outputFile.WriteString("\n")
			if err != nil {
				fmt.Println("Error writing newline:", err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through directory:", err)
		return
	}

	fmt.Println("Files concatenated successfully into 'combined.md'")
}

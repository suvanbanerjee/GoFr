package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Set the directory to read files from
	dirPath := "/Users/yuktha/Downloads/gogo" // Change this to your directory path

	// Create or open the output file
	outputFile, err := os.Create("output.txt") // Specify your desired output file
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Walk through the directory and process each file
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error walking the path:", err)
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Open each file for reading
		inputFile, err := os.Open(path)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return err
		}
		defer inputFile.Close()

		// Copy the contents of the current file to the output file
		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			fmt.Println("Error copying file contents:", err)
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through the directory:", err)
	}
}

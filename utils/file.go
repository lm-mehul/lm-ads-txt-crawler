package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func WriteStringArrayToFile(filename string, data []string) error {
	// Join the array elements into a single string
	content := strings.Join(data, "\n")
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return err
	}
	// Open or create the file for writing
	file, err := os.OpenFile(dir+"/resources/output/"+filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func ReadLinesFromFile(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

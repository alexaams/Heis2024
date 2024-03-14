package writefile

import (
	"fmt"
	"os"
	"path/filepath"
)

func checkFileExists(folderName, fileName string) error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working dir: %w", err)
	}

	folderPath := filepath.Join(workingDirectory, folderName)
	// Check if the folder exists, create it if not
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.Mkdir(folderPath, 0755); err != nil {
			return fmt.Errorf("creating folder '%s': %w", folderName, err)
		}
	}

	filePath := filepath.Join(folderPath, fileName)

	// Check if the file exists, create it if not
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("creating file '%s': %w", fileName, err)
		}
		defer file.Close()
	}

	return nil
}

func WriteToFile(folderName, fileName, msg string) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return
	}

	folderPath := filepath.Join(workingDirectory, folderName)
	filePath := filepath.Join(folderPath, fileName)

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644) // Ensure file is created if it doesn't exist
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, msg)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("Successfully wrote to the file:", filePath)
}

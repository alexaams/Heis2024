package backup

import (
	"fmt"
	"os"
	"path/filepath"
)

func checkFileExists(folderName, fileName string) error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getting Working dir: %w", err)
	}

	folderPath := filepath.Join(workingDirectory, folderName)
	// Check if the folder exists, create it if not
	if err := os.Mkdir(folderPath, 0755); err != nil { // 0755 permissions allow the owner to read/write/execute, others to read/execute
		return fmt.Errorf("creating folder '%s': %w", folderName, err)
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

func ReadFile(folderName, fileName string) {

}

package helpers

import (
	"log"
	"os"
	"path/filepath"
)

func FileExists(parentFolderDir string, fileName string) bool {
	// Create the full file path by joining the parent folder directory and the file name
	fullFilePath := filepath.Join(parentFolderDir, fileName)

	// Use os.Stat to get the file info
	_, err := os.Stat(fullFilePath)

	// If os.Stat returns an error and the error is of type os.ErrNotExist, the file does not exist
	if os.IsNotExist(err) {
		return false
	}

	// If there's no error, or an error other than os.ErrNotExist, we assume the file exists
	return err == nil
}

// Schreibt den Inhalt in eine Datei
// parameter:
// filePath: Pfad zur Datei
// content: Inhalt der in die Datei geschrieben werden soll
func WriteFile(filePath string, content string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Printf("Error writing to file %s: %v\n", filePath, err)
		return
	}
}

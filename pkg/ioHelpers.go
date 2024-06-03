package helpers

import (
	"log"
	"os"
	"path/filepath"
)

func FileExists(parentFolderDir string, fileName string) bool {
	fullFilePath := filepath.Join(parentFolderDir, fileName)
	_, err := os.Stat(fullFilePath)
	if os.IsNotExist(err) {
		return false
	}
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

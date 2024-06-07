package main

import (
	"io"
	"os"
	"testing"

	helpers "gihtub.com/VincentSchmid/whisper-transcription/pkg"
	subtitles "github.com/martinlindhe/subtitles"
)

func TestGetInvites(t *testing.T) {
	// read content from file to string
	content, _ := readFile("./output/transkribiert_CHECKUP_6497.srt")
	subs, _ := subtitles.NewFromSRT(content)
	newSubs := helpers.ConcatSubs(subs, 30)

	subString := newSubs.AsVTT()

	helpers.WriteFile("./output/concat.srt", subString)
}

func readFile(filename string) (string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

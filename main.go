package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	helpers "gihtub.com/VincentSchmid/whisper-transcription/pkg"
	subtitles "github.com/martinlindhe/subtitles"
	openai "github.com/sashabaranov/go-openai"
)

var (
	dataDir            string
	videoDir           string
	audioDir           string
	transcriptionDir   string
	gptResultDir       string
	transcribePrompt   string
	chatGptPrompt      string
	transcriptionSuffix string
)

func init() {
	dataDir = os.Getenv("DATA_DIR")
	videoDir = os.Getenv("VIDEO_DIR")
	audioDir = os.Getenv("AUDIO_DIR")
	transcriptionDir = os.Getenv("TRANSCRIPTION_DIR")
	gptResultDir = os.Getenv("GPT_RESULT_DIR")
	transcribePrompt = os.Getenv("TRANSCRIBE_PROMPT")
	chatGptPrompt = os.Getenv("CHAT_GPT_PROMPT")
	transcriptionSuffix = os.Getenv("TRANSCRIPTION_SUFFIX")

	if dataDir == "" || videoDir == "" || audioDir == "" || transcriptionDir == "" || gptResultDir == "" || transcribePrompt == "" || chatGptPrompt == "" || transcriptionSuffix == "" {
		log.Fatalf("One or more environment variables are not set")
	}
}

func main() {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		log.Fatalf("OPENAI_API_KEY is not set")
	}

	openaiClient := openai.NewClient(openaiKey)

	createFolders()

	err := convertVideoFilesToAudio(videoDir, audioDir)
	if err != nil {
		log.Println(err)
	}

	err = transcribeAudioFiles(openaiClient, audioDir, transcriptionSuffix, transcriptionDir)
	if err != nil {
		log.Println(err)
	}

	err = translateTranscriptions(openaiClient, transcriptionDir, transcriptionSuffix, gptResultDir)
	if err != nil {
		log.Println(err)
	}

	log.Println("Finished")
}

func createFolders() {
	log.Println("Creating data directories...")
	if err := os.MkdirAll(filepath.Join(videoDir), os.ModePerm); err != nil {
		log.Fatalf("Failed to create video directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(audioDir), os.ModePerm); err != nil {
		log.Fatalf("Failed to create audio directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(transcriptionDir), os.ModePerm); err != nil {
		log.Fatalf("Failed to create transcription directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(gptResultDir), os.ModePerm); err != nil {
		log.Fatalf("Failed to create GPT result directory: %v", err)
	}
}

func convertVideoFilesToAudio(inputDir string, outputDir string) error {
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			audioFilePath := filepath.Join(outputDir, baseName+".mp3")

			if !helpers.FileExists(outputDir, baseName+".mp3") {
				log.Printf("Converting %s to audio...\n", path)
				helpers.ConvertVideoToAudio(path, audioFilePath, 128)
			} else {
				log.Printf("File %s already exists\n", audioFilePath)
			}
		}

		return nil
	})
}

func transcribeAudioFiles(openaiClient *openai.Client, inputDir string, condesedSuffix string, outputDir string) error {
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			transcriptionPath := filepath.Join(outputDir, baseName+".vtt")
			condensedTranscriptionPath := filepath.Join(outputDir, baseName+condesedSuffix+".vtt")

			if !helpers.FileExists(outputDir, baseName+condesedSuffix+".vtt") {
				log.Printf("Transcribing audio file %s...\n", path)
				transcription := helpers.Transcribe(openaiClient, transcribePrompt, path)

				log.Printf("Writing transcription to %s...\n", transcriptionPath)
				helpers.WriteFile(transcriptionPath, transcription)

				log.Println("Condensing transcription timestamps...")
				subs, err := subtitles.NewFromVTT(transcription)
				if err != nil {
					return fmt.Errorf("error reading SRT file: %w", err)
				}

				condensedTranscript := helpers.ConcatSubs(subs, 30)

				log.Printf("Writing condensed transcription to %s...\n", condensedTranscriptionPath)
				helpers.WriteFile(condensedTranscriptionPath, condensedTranscript.AsSRT())
			} else {
				log.Printf("File %s already exists\n", condensedTranscriptionPath)
			}
		}

		return nil
	})
}

func translateTranscriptions(openaiClient *openai.Client, inputDir string, fileSuffix string, outputDir string) error {
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			resultPath := filepath.Join(outputDir, baseName+".txt")

			if strings.HasSuffix(baseName, fileSuffix) {

				if !helpers.FileExists(outputDir, baseName+".txt") {
					log.Printf("Reading Transcription from %s...\n", path)
					transcriptionBytes, err := os.ReadFile(path)

					if err != nil {
						return fmt.Errorf("error reading file %s: %w", path, err)
					}

					log.Println("Translating with chat gpt...")
					translatedResult, err := helpers.ChatGpt(openaiClient, chatGptPrompt, string(transcriptionBytes))

					if err != nil {
						return fmt.Errorf("error translating file %s: %w", path, err)
					}

					log.Printf("Writing chatGpt Response to %s...\n", resultPath)
					helpers.WriteFile(resultPath, translatedResult)

				} else {
					log.Printf("File %s already exists\n", resultPath)
				}
			}
		}

		return nil
	})
}

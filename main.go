package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	helpers "gihtub.com/VincentSchmid/whisper-transcription/pkg"
	"github.com/joho/godotenv"
	subtitles "github.com/martinlindhe/subtitles"
	openai "github.com/sashabaranov/go-openai"
)

var (
	openaiKey           string
	audioDir            string
	transcriptionDir    string
	gptResultDir        string
	transcribePrompt    string
	chatGptPrompt       string
	transcriptionSuffix = "_restructured"
)

func loadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is required", key)
	}
	return value
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	openaiKey = loadEnv("OPENAI_API_KEY")
	audioDir = loadEnv("AUDIO_DIR")
	transcriptionDir = loadEnv("TRANSCRIPTION_DIR")
	gptResultDir = loadEnv("GPT_RESULT_DIR")
	transcribePrompt = loadEnv("TRANSCRIBE_PROMPT")
	chatGptPrompt = loadEnv("CHAT_GPT_PROMPT")

	if _, err := os.Stat(audioDir); os.IsNotExist(err) {
		log.Printf("AUDIO_DIR %s does not exist", audioDir)
	}
}

func main() {
	openaiClient := openai.NewClient(openaiKey)

	createFolders()

	err := transcribeAudioFiles(openaiClient, audioDir, transcriptionSuffix, transcriptionDir)
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
	if err := os.MkdirAll(audioDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create audio directory: %v", err)
	}

	if err := os.MkdirAll(transcriptionDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create transcription directory: %v", err)
	}

	if err := os.MkdirAll(gptResultDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create GPT result directory: %v", err)
	}
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

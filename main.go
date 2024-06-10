package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	helpers "gihtub.com/VincentSchmid/whisper-transcription/pkg"
	"github.com/joho/godotenv"
	subtitles "github.com/martinlindhe/subtitles"
	openai "github.com/sashabaranov/go-openai"
)

var (
	devMode                 bool
	exeDir                  string
	openaiKey               string
	SubTitleTimeGranularity int
	audioDir                string
	transcriptionDir        string
	outputDir               string
	transcribePrompt        string
	chatGptPrompt           string
	transcriptionSuffix     = "_restructured"
)

func checkEnv(key string) bool {
	value := os.Getenv(key)
	return value == "true"
}

func loadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is required", key)
	}
	return value
}

func loadConfigFile() {
	envPath := filepath.Join(exeDir, "config.env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading config.env file")
	}
}

func init() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	exeDir = filepath.Dir(exePath)

	loadConfigFile()

	devMode = checkEnv("DEV_MODE")
	openaiKey = loadEnv("OPENAI_API_KEY")
	audioDir = loadEnv("AUDIO_DIR")
	transcriptionDir = loadEnv("TRANSCRIPTION_DIR")
	outputDir = loadEnv("OUTPUT_DIR")
	transcribePrompt = loadEnv("TRANSCRIBE_PROMPT")
	chatGptPrompt = loadEnv("CHAT_GPT_PROMPT")
	SubTitleTimeGranularity, err = strconv.Atoi(loadEnv("SUBTITLE_TIME_GRANULARITY"))
	if err != nil {
		log.Fatalf("Failed to parse SUBTITLE_TIME_GRANULARITY: %v", err)
	}

	// if is not absolute path, make it absolute
	if !filepath.IsAbs(audioDir) {
		audioDir = filepath.Join(exeDir, audioDir)
	}

	if !filepath.IsAbs(transcriptionDir) {
		transcriptionDir = filepath.Join(exeDir, transcriptionDir)
	}

	if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(exeDir, outputDir)
	}

	if _, err := os.Stat(audioDir); os.IsNotExist(err) {
		log.Printf("AUDIO_DIR %s does not exist", audioDir)
	}
}

func main() {
	if !devMode {
		f, err := os.OpenFile(filepath.Join(exeDir, "logs.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
	}

	openaiClient := openai.NewClient(openaiKey)

	createFolders()

	err := transcribeAudioFiles(openaiClient, audioDir, transcriptionSuffix, transcriptionDir)
	if err != nil {
		log.Println(err)
	}

	err = translateTranscriptions(openaiClient, transcriptionDir, transcriptionSuffix, outputDir)
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

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
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

				condensedTranscript := helpers.ConcatSubs(subs, SubTitleTimeGranularity)

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

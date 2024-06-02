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

const (
	dataDir = "./data"

	videoDir            = dataDir + "/1-video"
	audioDir            = dataDir + "/2-audio"
	transcriptionDir    = dataDir + "/3-transcription"
	gptResultDir        = dataDir + "/4-gpt_result"
	transcribePrompt    = "Es chömed Jugendlichi vor allem wo Risikoverhalte zeige, meischtens suizidali Jugendlichi, selte au Jugendlichi wo anderi bedrohe oder dor ihres Verhalte sich und anderi bedrohe. Do chas ned warte, die müend sofort igschätzt werde, wenn das kombiniert isch mit ere psychiatrische Uffälligkei. Aber es git au Störigsbilder wie Essstörige, schwerschti Depressione, manische Zustandsbilder oder Psychose, wo zwar au ned suizidal si müend, aber die chönnd au ned warte uf e reguläre Termin, au die chömed zu mir."
	chatGptPrompt       = "Übersetze die Text wort für wort is schwiizerdütsch. Lahne d'timestamps bestah."
	transcriptionSuffix = "_condensed"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		log.Fatalf("OPENAI_API_KEY is not set")
	}

	openaiClient := openai.NewClient(openaiKey)

	createFolders()

	err = convertVideoFilesToAudio(videoDir, audioDir)
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

func fileExists(parentFolderDir string, fileName string) bool {
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

func convertVideoFilesToAudio(inputDir string, outputDir string) error {
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			audioFilePath := filepath.Join(outputDir, baseName+".mp3")

			if !fileExists(outputDir, baseName+".mp3") {
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

			if !fileExists(outputDir, baseName+condesedSuffix+".vtt") {
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

				if !fileExists(outputDir, baseName+".txt") {
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

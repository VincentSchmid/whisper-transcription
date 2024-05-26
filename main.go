package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	subtitles "github.com/martinlindhe/subtitles"
	openai "github.com/sashabaranov/go-openai"
)

const (
	inputDir         = "./test/"
	outputDir        = "./output/"
	transcribePrompt = "Es chömed Jugendlichi vor allem wo Risikoverhalte zeige, meischtens suizidali Jugendlichi, selte au Jugendlichi wo anderi bedrohe oder dor ihres Verhalte sich und anderi bedrohe. Do chas ned warte, die müend sofort igschätzt werde, wenn das kombiniert isch mit ere psychiatrische Uffälligkei. Aber es git au Störigsbilder wie Essstörige, schwerschti Depressione, manische Zustandsbilder oder Psychose, wo zwar au ned suizidal si müend, aber die chönnd au ned warte uf e reguläre Termin, au die chömed zu mir."
	chatGptPrompt    = "Übersetze die Text wort für wort is schwiizerdütsch. Lahne d'timestamps bestah."
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	openaiKey := os.Getenv("OPENAI_API_KEY")

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	openaiClient := openai.NewClient(openaiKey)

	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		if !info.IsDir() {
			baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

			log.Printf("Transcribing %s\n", path)
			transcription := transcribe(openaiClient, transcribePrompt, path)

			log.Printf("Writing transcription\n")

			transcribeFilePath := filepath.Join(outputDir, "transkribiert_"+baseName+".srt")

			writeFile(transcribeFilePath, transcription)
			log.Printf("Transcription written to %s\n", transcribeFilePath)

			log.Printf("Condensing transcription timestamps\n")
			subs, err := subtitles.NewFromVTT(transcription)
			if err != nil {
				log.Printf("Error reading VTT file: %v\n", err)
				return err
			}

			condensedTranscript := concatSubs(subs, 30)

			log.Printf("Sending to GPT\n")
			gptResult, err := chatGpt(openaiClient, chatGptPrompt, condensedTranscript.AsSRT())
			if err != nil {
				log.Printf("ChatGPT error: %v\n", err)
				return err
			}

			log.Printf("Writing GPT result\n")

			gptResultFilePath := filepath.Join(outputDir, "schweizerdeutsch_"+baseName+".txt")

			writeFile(gptResultFilePath, gptResult)
			log.Printf("GPT result written to %s\n", gptResultFilePath)
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}
}

// Hier wird die Transkription mit der Zeitangabe zusammengefasst, sodass pro Zeitangabe weniger genau angezeigt wird
// parameter:
// subs: das Untertiel Objekt welche die Untertitel enthält
// lengthInSec: die Länge in Sekunden wie lange die Untertitel in einen Zeitblock zusammengefasst werden sollen
func concatSubs(subs subtitles.Subtitle, lengthInSec int) subtitles.Subtitle {
	var tmpCaption subtitles.Caption

	newSubs := subtitles.Subtitle{}
	startTime := subs.Captions[0].Start
	captionText := make([]string, 0)

	for i, caption := range subs.Captions {
		if caption.End.Sub(startTime).Seconds() < float64(lengthInSec) {
			captionText = append(captionText, caption.Text...)

		} else {
			tmpCaption = subtitles.Caption{
				Start: startTime,
				End:   caption.End,
				Text:  captionText,
			}

			newSubs.Captions = append(newSubs.Captions, tmpCaption)
			startTime = subs.Captions[i+1].Start
			captionText = make([]string, 0)
		}
	}
	return newSubs
}

// Whisper AI verwenden um das Audio file zu transkribieren
// parameter:
// openaiClient: das ist das openai Objekt welches den openai key hat
// prompt: Die Anweisung mit dem Beispiel Text für die AI
// filename: phad zur audio datei
func transcribe(openaiClient *openai.Client, prompt string, filename string) string {
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filename,
		Prompt:   prompt,
		Format:   openai.AudioResponseFormatVTT,
	}

	resp, err := openaiClient.CreateTranscription(ctx, req)
	if err != nil {
		log.Printf("Transcription error: %v\n", err)
		return ""
	}

	return resp.Text
}

// Schreibt den Inhalt in eine Datei
// parameter:
// filePath: Pfad zur Datei
// content: Inhalt der in die Datei geschrieben werden soll
func writeFile(filePath string, content string) {
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

// ChatGPT verwenden um den Text ins Schweizerdeutsch zu übersetzen
// parameter:
// openaiClient: das ist das openai Objekt welches den openai key hat
// systemPrompt: Die Anweisung für die AI was sie machen soll
// text: der Text der übersetzt werden soll (dies ist der Ihnalt der gesamten Datei)
func chatGpt(openaiClient *openai.Client, systemPrompt string, text string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		},
	}

	resp, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: messages,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

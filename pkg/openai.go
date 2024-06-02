package helpers

import (
	"context"
	"log"

	"github.com/sashabaranov/go-openai"
)

// Whisper AI verwenden um das Audio file zu transkribieren
// parameter:
// openaiClient: das ist das openai Objekt welches den openai key hat
// prompt: Die Anweisung mit dem Beispiel Text für die AI
// filename: phad zur audio datei
func Transcribe(openaiClient *openai.Client, prompt string, filename string) string {
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

// ChatGPT verwenden um den Text ins Schweizerdeutsch zu übersetzen
// parameter:
// openaiClient: das ist das openai Objekt welches den openai key hat
// systemPrompt: Die Anweisung für die AI was sie machen soll
// text: der Text der übersetzt werden soll (dies ist der Ihnalt der gesamten Datei)
func ChatGpt(openaiClient *openai.Client, systemPrompt string, text string) (string, error) {
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

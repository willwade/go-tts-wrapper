package main

import (
	"context"
	"log"

	"github.com/willwade/go-tts-wrapper/pkg/tts"
	_ "github.com/willwade/go-tts-wrapper/pkg/tts/providers/aws"
	_ "github.com/willwade/go-tts-wrapper/pkg/tts/providers/google"
)

func main() {
	config := tts.TTSConfig{
		APIKey:       "your-api-key",
		Region:       "us-west-1",
		LanguageCode: "en-US",
		VoiceID:      "en-US-Standard-A",
	}

	provider, err := tts.NewTTSProvider(tts.ProviderGoogle, config)
	if err != nil {
		log.Fatal(err)
	}

	err = provider.Speak(context.Background(), "Hello, world!")
	if err != nil {
		log.Fatal(err)
	}
}

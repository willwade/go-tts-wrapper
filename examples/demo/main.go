package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/willwade/go-tts-wrapper/pkg/tts"
	// Import only the providers you need
	_ "github.com/willwade/go-tts-wrapper/pkg/tts/providers/aws"    // AWS Polly
	_ "github.com/willwade/go-tts-wrapper/pkg/tts/providers/google" // Google Cloud TTS
)

func main() {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Load AWS credentials from environment
	config := tts.TTSConfig{
		APIKey:       os.Getenv("AWS_ACCESS_KEY"),
		Region:       os.Getenv("AWS_REGION"),
		LanguageCode: "en-US",
		VoiceID:      "Joanna",
		OutputFormat: "mp3",
	}

	// Initialize AWS Polly provider
	fmt.Println("Initializing AWS Polly provider...")
	provider, err := tts.NewTTSProvider(tts.ProviderAWS, config)
	if err != nil {
		log.Fatal(err)
	}

	// Basic text-to-speech
	fmt.Println("\nTesting basic text-to-speech...")
	if err := provider.Speak(ctx, "Hello! This is a test of the AWS Polly Text-to-Speech service."); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Demo completed successfully!")
} 
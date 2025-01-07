package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/willwade/go-tts-wrapper/pkg/tts"
	_ "github.com/willwade/go-tts-wrapper/pkg/tts/providers/aws"
	_ "github.com/willwade/go-tts-wrapper/pkg/tts/providers/google"
	_ "github.com/willwade/go-tts-wrapper/pkg/tts/providers/ibm"
)

func main() {
	if err := runDemo(); err != nil {
		log.Fatal(err)
	}
}

func runDemo() error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Demo configuration
	config := tts.TTSConfig{
		APIKey:       os.Getenv("AWS_ACCESS_KEY"),    // For AWS Polly
		Region:       "us-west-1",
		LanguageCode: "en-US",
		VoiceID:      "Joanna",                       // AWS Polly voice
		OutputFormat: "mp3",
	}

	// Initialize provider (AWS Polly in this example)
	fmt.Println("Initializing AWS Polly provider...")
	provider, err := tts.NewTTSProvider(tts.ProviderAWS, config)
	if err != nil {
		return fmt.Errorf("failed to initialize provider: %w", err)
	}

	// List available voices
	fmt.Println("\nListing available voices...")
	voices, err := provider.GetVoices(ctx)
	if err != nil {
		return fmt.Errorf("failed to get voices: %w", err)
	}

	fmt.Println("Available voices:")
	for i, voice := range voices {
		if i >= 5 {
			fmt.Println("... (more voices available)")
			break
		}
		fmt.Printf("- %s (%s, %s)\n", voice.Name, voice.Language, voice.Gender)
	}

	// List audio devices
	fmt.Println("\nListing audio devices...")
	devices, err := tts.ListAudioDevices()
	if err != nil {
		return fmt.Errorf("failed to list audio devices: %w", err)
	}

	fmt.Println("Available audio devices:")
	for _, device := range devices {
		fmt.Printf("- %s (Default: %v)\n", device.Name, device.IsDefault)
		if device.IsDefault {
			if err := provider.SetOutputDevice(device.ID); err != nil {
				return fmt.Errorf("failed to set output device: %w", err)
			}
		}
	}

	// Basic text-to-speech
	fmt.Println("\nDemonstrating basic text-to-speech...")
	text := "Hello! This is a demonstration of the Go TTS Wrapper library."
	if err := provider.Speak(ctx, text); err != nil {
		return fmt.Errorf("failed to speak text: %w", err)
	}

	// Demonstrate speech properties
	fmt.Println("\nDemonstrating different speech properties...")
	
	// Faster speech
	fmt.Println("Speaking faster...")
	if err := provider.SetProperty("rate", 1.5); err != nil {
		return fmt.Errorf("failed to set rate: %w", err)
	}
	if err := provider.Speak(ctx, "This is what faster speech sounds like."); err != nil {
		return fmt.Errorf("failed to speak faster: %w", err)
	}

	// Reset rate
	if err := provider.SetProperty("rate", 1.0); err != nil {
		return fmt.Errorf("failed to reset rate: %w", err)
	}

	// SSML demonstration
	fmt.Println("\nDemonstrating SSML capabilities...")
	ssml := `<speak>
		Welcome to <break time="0.5s"/> SSML demonstration.
		<prosody rate="slow">This is spoken slowly,</prosody>
		<prosody pitch="+2st">and this is spoken in a higher pitch.</prosody>
		<break time="1s"/>
		<emphasis level="strong">Thank you for listening!</emphasis>
	</speak>`

	if err := provider.SpeakSSML(ctx, ssml); err != nil {
		return fmt.Errorf("failed to speak SSML: %w", err)
	}

	// Demonstrate audio controls
	fmt.Println("\nDemonstrating audio controls...")
	longText := "This is a longer piece of text that we'll use to demonstrate audio controls. " +
		"While this is playing, we'll pause, resume, and then stop the playback."
	
	// Start speaking in a goroutine
	go func() {
		if err := provider.Speak(ctx, longText); err != nil {
			fmt.Printf("Error during playback: %v\n", err)
		}
	}()

	// Wait a bit then pause
	time.Sleep(2 * time.Second)
	fmt.Println("Pausing...")
	if err := provider.PauseAudio(); err != nil {
		return fmt.Errorf("failed to pause: %w", err)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("Resuming...")
	if err := provider.ResumeAudio(); err != nil {
		return fmt.Errorf("failed to resume: %w", err)
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Stopping...")
	if err := provider.StopAudio(); err != nil {
		return fmt.Errorf("failed to stop: %w", err)
	}

	// File output demonstration
	fmt.Println("\nDemonstrating file output...")
	if err := provider.SynthToFile(ctx, "This text will be saved to a file.", "output.mp3"); err != nil {
		return fmt.Errorf("failed to synthesize to file: %w", err)
	}
	fmt.Println("Audio has been saved to 'output.mp3'")

	fmt.Println("\nDemo completed successfully!")
	return nil
} 
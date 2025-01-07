# Go Text-to-Speech Wrapper

A flexible Go library that provides a unified interface for multiple text-to-speech (TTS) providers including AWS Polly, Google Cloud TTS, Microsoft Azure, IBM Watson, ElevenLabs, eSpeak-NG, and Sherpa-ONNX.

## Features

- Multiple TTS provider support
- Unified interface for all providers
- Real-time audio playback
- SSML support (where available)
- Voice selection
- Audio control (pause/resume/stop)
- Configurable speech properties (rate, pitch, volume)
- Audio device selection
- Caching support

## Installation

```bash
go get github.com/willwade/go-tts-wrapper
```

## Basic Usage

```go
package main

import (
    "context"
    "log"

    "github.com/willwade/go-tts-wrapper/pkg/tts"
    _ "github.com/willwade/go-tts-wrapper/pkg/tts/providers/aws"    // Import AWS provider
    _ "github.com/willwade/go-tts-wrapper/pkg/tts/providers/google" // Import Google provider
)

func main() {
    // Create a configuration for your chosen provider
    config := tts.TTSConfig{
        APIKey:       "your-api-key",
        Region:       "us-west-1",    // Required for AWS, Azure
        LanguageCode: "en-US",
        VoiceID:      "en-US-Standard-A",
    }

    // Create a provider (e.g., Google Cloud TTS)
    provider, err := tts.NewTTSProvider(tts.ProviderGoogle, config)
    if err != nil {
        log.Fatal(err)
    }

    // Speak some text
    err = provider.Speak(context.Background(), "Hello, world!")
    if err != nil {
        log.Fatal(err)
    }
}
```

## Configuration

You can configure the TTS wrapper using either code or a JSON configuration file:

```go
// Load configuration from file
config, err := config.LoadConfig("config.json")
if err != nil {
    log.Fatal(err)
}
```

Example config.json:
```json
{
    "providers": {
        "aws": {
            "apiKey": "YOUR_AWS_ACCESS_KEY",
            "region": "us-west-1",
            "voiceId": "Joanna"
        },
        "google": {
            "languageCode": "en-US",
            "voiceId": "en-US-Standard-A"
        }
    },
    "audio": {
        "rate": 1.0,
        "pitch": 1.0,
        "volume": 1.0,
        "deviceId": ""
    },
    "cache": {
        "enabled": true,
        "directory": "./cache",
        "maxSize": 1073741824,
        "ttl": 86400,
        "filePattern": "tts-{hash}.mp3"
    }
}
```

## Provider-Specific Configuration

### AWS Polly
```go
config := tts.TTSConfig{
    APIKey:  "YOUR_AWS_ACCESS_KEY",    // AWS credentials should be configured
    Region:  "us-west-1",
    VoiceID: "Joanna",
}
provider, _ := tts.NewTTSProvider(tts.ProviderAWS, config)
```

### Google Cloud TTS
```go
config := tts.TTSConfig{
    // Use GOOGLE_APPLICATION_CREDENTIALS environment variable
    LanguageCode: "en-US",
    VoiceID:     "en-US-Standard-A",
}
provider, _ := tts.NewTTSProvider(tts.ProviderGoogle, config)
```

### Microsoft Azure
```go
config := tts.TTSConfig{
    APIKey:  "YOUR_AZURE_KEY",
    Region:  "eastus",
    VoiceID: "en-US-JennyNeural",
}
provider, _ := tts.NewTTSProvider(tts.ProviderMicrosoft, config)
```

### IBM Watson
```go
config := tts.TTSConfig{
    APIKey:  "YOUR_IBM_KEY",
    Region:  "us-south",
    VoiceID: "en-US_AllisonV3Voice",
}
provider, _ := tts.NewTTSProvider(tts.ProviderIBM, config)
```

### ElevenLabs
```go
config := tts.TTSConfig{
    APIKey:  "YOUR_ELEVENLABS_KEY",
    VoiceID: "voice-id",
}
provider, _ := tts.NewTTSProvider(tts.ProviderElevenLabs, config)
```

### eSpeak-NG (Local)
```go
config := tts.TTSConfig{
    VoiceID: "en",  // Language/voice code
}
provider, _ := tts.NewTTSProvider(tts.ProviderESpeak, config)
```

### Sherpa-ONNX (Local)
```go
config := tts.TTSConfig{
    Engine: "/path/to/model.onnx",  // Path to VITS model
}
provider, _ := tts.NewTTSProvider(tts.ProviderSherpaONNX, config)
```

## Advanced Usage

### Using SSML
```go
ssml := `<speak>
    Hello <break time="1s"/> World!
    <prosody rate="slow">This is slower</prosody>
</speak>`

err := provider.SpeakSSML(context.Background(), ssml)
```

### Listing Available Voices
```go
voices, err := provider.GetVoices(context.Background())
if err != nil {
    log.Fatal(err)
}

for _, voice := range voices {
    fmt.Printf("ID: %s, Name: %s, Language: %s\n", 
        voice.ID, voice.Name, voice.Language)
}
```

### Audio Device Selection
```go
// List available audio devices
devices, err := tts.ListAudioDevices()
if err != nil {
    log.Fatal(err)
}

// Set output device
provider.SetOutputDevice(devices[0].ID)
```

### Controlling Audio Playback
```go
// Start speaking
provider.Speak(context.Background(), "This is a long text...")

// Pause playback
provider.PauseAudio()

// Resume playback
provider.ResumeAudio()

// Stop playback
provider.StopAudio()
```

### Adjusting Speech Properties
```go
provider.SetProperty("rate", 1.5)    // Speed up speech
provider.SetProperty("pitch", 0.8)   // Lower pitch
provider.SetProperty("volume", 1.2)  // Increase volume
```

## Dependencies

- PortAudio for audio playback
- Provider-specific SDKs (automatically managed through Go modules)
- For eSpeak-NG: `espeak-ng` binary must be installed on the system
- For Sherpa-ONNX: Requires ONNX runtime and model files

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

package tts

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// MicrosoftProvider implements TTSProvider for Azure Cognitive Services
type MicrosoftProvider struct {
	*BaseProvider
	config      *speech.SpeechConfig
	audioPlayer *AudioPlayer
}

// NewMicrosoftProvider creates a new Azure TTS provider
func NewMicrosoftProvider(cfg TTSConfig) (*MicrosoftProvider, error) {
	speechConfig, err := speech.NewSpeechConfigFromSubscription(cfg.APIKey, cfg.Region)
	if err != nil {
		return nil, err
	}

	audioPlayer, err := NewAudioPlayer()
	if err != nil {
		return nil, err
	}

	return &MicrosoftProvider{
		BaseProvider: NewBaseProvider(cfg),
		config:       speechConfig,
		audioPlayer:  audioPlayer,
	}, nil
}

func (p *MicrosoftProvider) synthesize(ctx context.Context, text string, isSSML bool) ([]byte, error) {
	synthesizer, err := speech.NewSpeechSynthesizerFromConfig(p.config, nil)
	if err != nil {
		return nil, err
	}
	defer synthesizer.Close()

	var result *speech.SpeechSynthesisResult
	if isSSML {
		result, err = synthesizer.SpeakSsmlAsync(text).Get()
	} else {
		result, err = synthesizer.SpeakTextAsync(text).Get()
	}
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if result.Reason != common.SynthesizingAudioCompleted {
		return nil, fmt.Errorf("synthesis failed: %v", result.Reason)
	}

	return result.AudioData, nil
}

func (p *MicrosoftProvider) Speak(ctx context.Context, text string) error {
	audioData, err := p.synthesize(ctx, text, false)
	if err != nil {
		return err
	}
	return p.audioPlayer.PlayMP3Stream(io.NopCloser(bytes.NewReader(audioData)))
}

func (p *MicrosoftProvider) SpeakSSML(ctx context.Context, ssml string) error {
	if err := p.ValidateSSML(ssml); err != nil {
		return err
	}
	audioData, err := p.synthesize(ctx, ssml, true)
	if err != nil {
		return err
	}
	return p.audioPlayer.PlayMP3Stream(io.NopCloser(bytes.NewReader(audioData)))
}

func (p *MicrosoftProvider) GetVoices(ctx context.Context) ([]Voice, error) {
	synthesizer, err := speech.NewSpeechSynthesizerFromConfig(p.config, nil)
	if err != nil {
		return nil, err
	}
	defer synthesizer.Close()

	voicesList, err := synthesizer.GetVoicesAsync("").Get()
	if err != nil {
		return nil, err
	}

	voices := make([]Voice, 0, len(voicesList.Voices))
	for _, v := range voicesList.Voices {
		voices = append(voices, Voice{
			ID:          v.ShortName,
			Name:        v.DisplayName,
			Language:    v.Locale,
			Gender:      v.Gender,
			Provider:    "Microsoft",
			NativeVoice: v,
		})
	}
	return voices, nil
}

// Audio control methods
func (p *MicrosoftProvider) PauseAudio() error {
	return p.audioPlayer.Pause()
}

func (p *MicrosoftProvider) ResumeAudio() error {
	return p.audioPlayer.Resume()
}

func (p *MicrosoftProvider) StopAudio() error {
	return p.audioPlayer.Stop()
}

func (p *MicrosoftProvider) CheckCredentials(ctx context.Context) bool {
	synthesizer, err := speech.NewSpeechSynthesizerFromConfig(p.config, nil)
	if err != nil {
		return false
	}
	defer synthesizer.Close()

	_, err = synthesizer.GetVoicesAsync("").Get()
	return err == nil
}

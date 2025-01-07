package tts

import (
	"context"
	"io"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/watson-developer-cloud/go-sdk/texttospeechv1"
	"github.com/willwade/go-tts-wrapper/pkg/tts"
)

// IBMProvider implements TTSProvider for IBM Watson
type IBMProvider struct {
	*tts.BaseProvider
	client      *texttospeechv1.TextToSpeechV1
	audioPlayer *tts.AudioPlayer
}

// NewIBMProvider creates a new IBM Watson TTS provider
func NewIBMProvider(cfg tts.TTSConfig) (*IBMProvider, error) {
	authenticator := &core.IamAuthenticator{
		ApiKey: cfg.APIKey,
	}

	service, err := texttospeechv1.NewTextToSpeechV1(&texttospeechv1.TextToSpeechV1Options{
		Authenticator: authenticator,
		URL:           "https://api.us-south.text-to-speech.watson.cloud.ibm.com", // Default region
	})
	if err != nil {
		return nil, err
	}

	audioPlayer, err := NewAudioPlayer()
	if err != nil {
		return nil, err
	}

	return &IBMProvider{
		BaseProvider: NewBaseProvider(cfg),
		client:       service,
		audioPlayer:  audioPlayer,
	}, nil
}

func (p *IBMProvider) synthesize(ctx context.Context, text string, isSSML bool) (io.ReadCloser, error) {
	var synthesizeOptions *texttospeechv1.SynthesizeOptions

	if isSSML {
		synthesizeOptions = service.NewSynthesizeOptions(text).
			SetAccept("audio/mp3").
			SetVoice(p.config.VoiceID).
			SetTextType("ssml")
	} else {
		synthesizeOptions = service.NewSynthesizeOptions(text).
			SetAccept("audio/mp3").
			SetVoice(p.config.VoiceID)
	}

	result, _, err := p.client.Synthesize(synthesizeOptions)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *IBMProvider) Speak(ctx context.Context, text string) error {
	audio, err := p.synthesize(ctx, text, false)
	if err != nil {
		return err
	}
	defer audio.Close()

	return p.audioPlayer.PlayMP3Stream(audio)
}

func (p *IBMProvider) SpeakSSML(ctx context.Context, ssml string) error {
	if err := p.ValidateSSML(ssml); err != nil {
		return err
	}
	audio, err := p.synthesize(ctx, ssml, true)
	if err != nil {
		return err
	}
	defer audio.Close()

	return p.audioPlayer.PlayMP3Stream(audio)
}

func (p *IBMProvider) GetVoices(ctx context.Context) ([]tts.Voice, error) {
	result, _, err := p.client.ListVoices(service.NewListVoicesOptions())
	if err != nil {
		return nil, err
	}

	voices := make([]tts.Voice, 0, len(result.Voices))
	for _, v := range result.Voices {
		voices = append(voices, tts.Voice{
			ID:          *v.Name,
			Name:        *v.Name,
			Language:    *v.Language,
			Gender:      *v.Gender,
			Provider:    "IBM",
			NativeVoice: v,
		})
	}
	return voices, nil
}

// Audio control methods
func (p *IBMProvider) PauseAudio() error {
	return p.audioPlayer.Pause()
}

func (p *IBMProvider) ResumeAudio() error {
	return p.audioPlayer.Resume()
}

func (p *IBMProvider) StopAudio() error {
	return p.audioPlayer.Stop()
}

func (p *IBMProvider) CheckCredentials(ctx context.Context) bool {
	_, _, err := p.client.ListVoices(service.NewListVoicesOptions())
	return err == nil
}

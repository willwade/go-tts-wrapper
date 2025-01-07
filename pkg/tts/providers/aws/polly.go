package tts

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
	tts "github.com/willwade/go-tts-wrapper"
)

type PollyProvider struct {
	*tts.BaseProvider
	client      *polly.Client
	audioPlayer *tts.AudioPlayer
}

func NewPollyProvider(cfg tts.TTSConfig) (*PollyProvider, error) {
	// Load AWS configuration
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	// Create Polly client
	client := polly.NewFromConfig(awsCfg)

	audioPlayer, err := tts.NewAudioPlayer()
	if err != nil {
		return nil, err
	}

	return &PollyProvider{
		BaseProvider: tts.NewBaseProvider(cfg),
		client:      client,
		audioPlayer: audioPlayer,
	}, nil
}

func (p *PollyProvider) synthesize(ctx context.Context, text string, isSSML bool) ([]byte, error) {
	inputType := types.TextTypeText
	if isSSML {
		inputType = types.TextTypeSsml
	}

	input := &polly.SynthesizeSpeechInput{
		Engine:       types.EngineNeural,
		OutputFormat: types.OutputFormatMp3,
		Text:         &text,
		TextType:     inputType,
		VoiceId:      types.VoiceId(p.config.VoiceID),
	}

	resp, err := p.client.SynthesizeSpeech(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to synthesize speech: %w", err)
	}
	defer resp.AudioStream.Close()

	return io.ReadAll(resp.AudioStream)
}

func (p *PollyProvider) Speak(ctx context.Context, text string) error {
	audioData, err := p.synthesize(ctx, text, false)
	if err != nil {
		return err
	}
	return p.audioPlayer.PlayMP3Stream(bytes.NewReader(audioData))
}

func (p *PollyProvider) SpeakSSML(ctx context.Context, ssml string) error {
	if err := p.ValidateSSML(ssml); err != nil {
		return err
	}
	audioData, err := p.synthesize(ctx, ssml, true)
	if err != nil {
		return err
	}
	return p.audioPlayer.PlayMP3Stream(bytes.NewReader(audioData))
}

func (p *PollyProvider) GetVoices(ctx context.Context) ([]tts.Voice, error) {
	input := &polly.DescribeVoicesInput{
		Engine: types.EngineNeural,
	}

	resp, err := p.client.DescribeVoices(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get voices: %w", err)
	}

	voices := make([]tts.Voice, 0, len(resp.Voices))
	for _, v := range resp.Voices {
		voices = append(voices, tts.Voice{
			ID:          string(v.Id),
			Name:        *v.Name,
			Language:    *v.LanguageCode,
			Gender:      string(v.Gender),
			Provider:    "AWS Polly",
			NativeVoice: v,
		})
	}

	return voices, nil
}

// Audio control methods
func (p *PollyProvider) PauseAudio() error {
	return p.audioPlayer.Pause()
}

func (p *PollyProvider) ResumeAudio() error {
	return p.audioPlayer.Resume()
}

func (p *PollyProvider) StopAudio() error {
	return p.audioPlayer.Stop()
}

func (p *PollyProvider) CheckCredentials(ctx context.Context) bool {
	_, err := p.client.DescribeVoices(ctx, &polly.DescribeVoicesInput{})
	return err == nil
}

func init() {
	tts.RegisterProvider(tts.ProviderAWS, func(cfg tts.TTSConfig) (tts.TTSProvider, error) {
		return NewPollyProvider(cfg)
	})
} 
package tts

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ESpeakProvider implements TTSProvider for eSpeak-NG
type ESpeakProvider struct {
	*BaseProvider
	audioPlayer *AudioPlayer
}

// NewESpeakProvider creates a new eSpeak-NG provider
func NewESpeakProvider(cfg TTSConfig) (*ESpeakProvider, error) {
	// Check if espeak-ng is installed
	if _, err := exec.LookPath("espeak-ng"); err != nil {
		return nil, fmt.Errorf("espeak-ng not found: %w", err)
	}

	audioPlayer, err := NewAudioPlayer()
	if err != nil {
		return nil, err
	}

	return &ESpeakProvider{
		BaseProvider: NewBaseProvider(cfg),
		audioPlayer:  audioPlayer,
	}, nil
}

func (p *ESpeakProvider) synthesize(ctx context.Context, text string, isSSML bool) ([]byte, error) {
	args := []string{
		"--stdout",
		"--voice=" + p.config.VoiceID,
		"--rate=" + strconv.Itoa(int(p.audioConfig.Rate*175)),     // Default rate is 175 words per minute
		"--pitch=" + strconv.Itoa(int(p.audioConfig.Pitch*50)),    // Range 0-100
		"--volume=" + strconv.Itoa(int(p.audioConfig.Volume*100)), // Range 0-200
	}

	if isSSML {
		args = append(args, "-m") // Enable SSML/markup
	}

	cmd := exec.CommandContext(ctx, "espeak-ng", append(args, text)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("espeak-ng failed: %w: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

func (p *ESpeakProvider) Speak(ctx context.Context, text string) error {
	audioData, err := p.synthesize(ctx, text, false)
	if err != nil {
		return err
	}
	return p.audioPlayer.PlayMP3Stream(bytes.NewReader(audioData))
}

func (p *ESpeakProvider) SpeakSSML(ctx context.Context, ssml string) error {
	if err := p.ValidateSSML(ssml); err != nil {
		return err
	}
	audioData, err := p.synthesize(ctx, ssml, true)
	if err != nil {
		return err
	}
	return p.audioPlayer.PlayMP3Stream(bytes.NewReader(audioData))
}

func (p *ESpeakProvider) GetVoices(ctx context.Context) ([]Voice, error) {
	cmd := exec.CommandContext(ctx, "espeak-ng", "--voices")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get voices: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	voices := make([]Voice, 0, len(lines)-1)

	// Skip header line
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		voices = append(voices, Voice{
			ID:       fields[1],
			Language: fields[2],
			Gender:   fields[3],
			Name:     strings.Join(fields[4:], " "),
			Provider: "eSpeak-NG",
		})
	}

	return voices, nil
}

// Audio control methods
func (p *ESpeakProvider) PauseAudio() error {
	return p.audioPlayer.Pause()
}

func (p *ESpeakProvider) ResumeAudio() error {
	return p.audioPlayer.Resume()
}

func (p *ESpeakProvider) StopAudio() error {
	return p.audioPlayer.Stop()
}

func (p *ESpeakProvider) SetOutputDevice(deviceID string) error {
	return ErrNotImplemented
}

func (p *ESpeakProvider) CheckCredentials(ctx context.Context) bool {
	// eSpeak-NG is a local binary, so just check if it's available
	_, err := exec.LookPath("espeak-ng")
	return err == nil
}

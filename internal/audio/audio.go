package pkg

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/hajimehoshi/go-mp3"
)

// AudioPlayer handles audio playback using PortAudio
type AudioPlayer struct {
	stream     *portaudio.Stream
	buffer     []float32
	playing    bool
	pauseLock  sync.Mutex
	paused     bool
	done       chan struct{}
	sampleRate float64
}

// NewAudioPlayer creates a new audio player instance
func NewAudioPlayer() (*AudioPlayer, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize PortAudio: %w", err)
	}

	return &AudioPlayer{
		sampleRate: 44100,
		done:       make(chan struct{}),
	}, nil
}

// PlayMP3Stream plays MP3 audio data from an io.Reader
func (ap *AudioPlayer) PlayMP3Stream(r io.Reader) error {
	// Read all audio data into memory
	audioData, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read audio data: %w", err)
	}

	// Convert MP3 to PCM
	pcmData, sampleRate, err := mp3ToPCM(bytes.NewReader(audioData))
	if err != nil {
		return fmt.Errorf("failed to decode MP3: %w", err)
	}

	ap.pauseLock.Lock()
	ap.buffer = pcmData
	ap.playing = true
	ap.paused = false
	ap.sampleRate = sampleRate
	ap.pauseLock.Unlock()

	// Reset completion channel
	select {
	case <-ap.done:
	default:
	}
	ap.done = make(chan struct{})

	// Initialize the PortAudio stream
	stream, err := portaudio.OpenDefaultStream(0, 1, ap.sampleRate,
		len(pcmData), ap.buffer)
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %w", err)
	}
	ap.stream = stream

	if err := stream.Start(); err != nil {
		stream.Close()
		return fmt.Errorf("failed to start audio stream: %w", err)
	}

	// Write audio data to the stream in chunks
	go func() {
		defer close(ap.done)
		if err := stream.Write(); err != nil {
			fmt.Printf("Error writing to audio stream: %v\n", err)
		}
		ap.playing = false
	}()

	return nil
}

// WaitForCompletion blocks until playback is complete or context is cancelled
func (ap *AudioPlayer) WaitForCompletion() error {
	if ap.done == nil {
		return fmt.Errorf("no active playback")
	}
	<-ap.done
	return nil
}

// IsPlaying returns true if audio is currently playing
func (ap *AudioPlayer) IsPlaying() bool {
	ap.pauseLock.Lock()
	defer ap.pauseLock.Unlock()
	return ap.playing && !ap.paused
}

// Pause pauses audio playback
func (ap *AudioPlayer) Pause() error {
	ap.pauseLock.Lock()
	defer ap.pauseLock.Unlock()

	if ap.stream == nil || !ap.playing {
		return fmt.Errorf("no active audio playback")
	}

	if err := ap.stream.Stop(); err != nil {
		return fmt.Errorf("failed to pause audio: %w", err)
	}
	ap.paused = true
	return nil
}

// Resume resumes audio playback
func (ap *AudioPlayer) Resume() error {
	ap.pauseLock.Lock()
	defer ap.pauseLock.Unlock()

	if ap.stream == nil || !ap.playing || !ap.paused {
		return fmt.Errorf("no paused audio playback")
	}

	if err := ap.stream.Start(); err != nil {
		return fmt.Errorf("failed to resume audio: %w", err)
	}
	ap.paused = false
	return nil
}

// Stop stops audio playback and cleans up resources
func (ap *AudioPlayer) Stop() error {
	ap.pauseLock.Lock()
	defer ap.pauseLock.Unlock()

	if ap.stream == nil {
		return nil
	}

	if err := ap.stream.Stop(); err != nil {
		return fmt.Errorf("failed to stop audio: %w", err)
	}

	if err := ap.stream.Close(); err != nil {
		return fmt.Errorf("failed to close audio stream: %w", err)
	}

	ap.stream = nil
	ap.playing = false
	ap.paused = false
	close(ap.done)
	return nil
}

// Close cleans up the AudioPlayer resources
func (ap *AudioPlayer) Close() error {
	if err := ap.Stop(); err != nil {
		return err
	}
	return portaudio.Terminate()
}

// mp3ToPCM converts MP3 data to PCM format using go-mp3
func mp3ToPCM(r io.Reader) ([]float32, float64, error) {
	decoder, err := mp3.NewDecoder(r)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create MP3 decoder: %w", err)
	}

	// Create a buffer for the PCM data
	var pcmData []float32
	buffer := make([]byte, 4096)

	for {
		n, err := decoder.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("error reading MP3 data: %w", err)
		}

		// Convert each pair of bytes (16-bit samples) to float32
		for i := 0; i < n; i += 2 {
			if i+1 >= n {
				break
			}
			// Convert 16-bit PCM to float32 (-1.0 to 1.0)
			sample := float32(int16(buffer[i])|int16(buffer[i+1])<<8) / 32768.0
			pcmData = append(pcmData, sample)
		}
	}

	if len(pcmData) == 0 {
		return nil, 0, fmt.Errorf("no audio data decoded")
	}

	return pcmData, float64(decoder.SampleRate()), nil
}

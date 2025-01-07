// pkg/tts/audio/devices.go
package tts

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

type AudioDevice struct {
    ID          string
    Name        string
    IsDefault   bool
    SampleRates []int
}

func ListAudioDevices() ([]AudioDevice, error) {
    err := portaudio.Initialize()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize PortAudio: %w", err)
    }
    defer portaudio.Terminate()

    devices, err := portaudio.Devices()
    if err != nil {
        return nil, fmt.Errorf("failed to get devices: %w", err)
    }

    audioDevices := make([]AudioDevice, 0, len(devices))
    for _, dev := range devices {
        if dev.MaxOutputChannels > 0 {
            audioDevices = append(audioDevices, AudioDevice{
                ID:          dev.Name,
                Name:        dev.Name,
                IsDefault:   dev.Name == portaudio.DefaultOutputDevice.Name,
                SampleRates: []int{44100, 48000}, // Common sample rates
            })
        }
    }
    return audioDevices, nil
}

func SetDefaultDevice(deviceID string) error {
    devices, err := ListAudioDevices()
    if err != nil {
        return err
    }

    for _, dev := range devices {
        if dev.ID == deviceID {
            return nil // Device exists
        }
    }
    return fmt.Errorf("device not found: %s", deviceID)
}
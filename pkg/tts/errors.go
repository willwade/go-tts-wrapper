package tts

import "fmt"

type TTSError struct {
    Provider string
    Code     string
    Message  string
    Err      error
}

func (e *TTSError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%s] %s: %v", e.Provider, e.Message, e.Err)
    }
    return fmt.Sprintf("[%s] %s", e.Provider, e.Message)
}
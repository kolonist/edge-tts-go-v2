// Package edgetts contains structs and functions to generate speech from text
package edgetts

import (
	"fmt"
	"iter"

	"github.com/kolonist/edgetts/internal/tts"
)

// Args contains parametsrs for TTS
type Args = tts.Args

// OutputFormat represents sound data output format
type OutputFormat = tts.OutputFormat

// Sound data output formats possible values
const (
	// mp3 24khz, 48k bitrate (default)
	OutputFormatMp3 = tts.OutputFormatMp3

	// webm 24khz, 16bit, 24k bitrate
	OutputFormatWebm = tts.OutputFormatWebm

	// ogg 24khz, 16bit
	OutputFormatOgg = tts.OutputFormatOgg

	// raw PCM 22050 hz, 16bit
	OutputFormatRaw22050 = tts.OutputFormatRaw22050

	// raw PCM 44100 hz, 16bit
	OutputFormatRaw44100 = tts.OutputFormatRaw44100
)

// SpeechMetadata contains time of word start and its pronunciation duration
type SpeechMetadata = tts.SpeechMetadata

// EdgeTTS used to generate speech from text
type EdgeTTS struct {
	text     string
	args     Args
	ready    bool
	metadata []SpeechMetadata
}

// New creates EdgeTTS struct with arguments to generate speech.
//
// Parameters:
//
//	args - parameters of speech generation like voice, rate, volume
//
// Returns:
//
//	New EdgeTTS struct
func New(args Args) *EdgeTTS {
	return &EdgeTTS{
		text:     "",
		args:     args,
		ready:    false,
		metadata: nil,
	}
}

// Speak define text you need to convert to speech
//
// Parameters:
//
//	text - text to speak
//
// Returns:
//
//	Current EdgeTTS struct to use in chained calls
func (etts *EdgeTTS) Speak(text string) *EdgeTTS {
	etts.text = text
	return etts
}

// Speak define text you need to convert to speech
//
// Parameters:
//
//	text - text to speak
//	voice - voice os speaker to read text
//
// Returns:
//
//	Current EdgeTTS struct to use in chained calls
func (etts *EdgeTTS) SpeakWithVoice(text string, voice string) *EdgeTTS {
	etts.text = text
	etts.args.Voice = voice

	return etts
}

// SaveToFile generate speach and save it to file.
//
// Parameters:
//
//	filename - full path to file. Should be write accessible
//	format - format of sound data to write to file. Use one of OutputFormat* constants
//
// Returns:
//
//	error if file was not written for some reason
func (etts *EdgeTTS) SaveToFile(filename string, format OutputFormat) error {
	return nil
}

// GetSoundBytesIterator generate speech and return it in iterator with small byte buffers as they come from server.
//
// Parameters:
//
//	format - format of sound data to write to file. Use one of OutputFormat* constants
//
// Returns:
//
//	error if file was not written for some reason
func (etts *EdgeTTS) GetSoundBytesIterator(format OutputFormat) iter.Seq[[]byte] {
	return nil
}

// GetSoundBytes generate speech and returns it as bytes.
//
// Parameters:
//
//	format - format of sound data to write to file. Use one of OutputFormat* constants
//
// Returns:
//
//   - buffer containing sound data of defined format
//   - error if file was not written for some reason
func (etts *EdgeTTS) GetSoundBytes(format OutputFormat) ([]byte, error) {
	return nil, nil
}

// GetMetadata gets metadata of generated speech
//
// Returns:
//
//   - Slice with SpeechMetadata structs containing timings of each word in text
//   - Error
func (etts *EdgeTTS) GetMetadata() ([]SpeechMetadata, error) {
	if etts.ready {
		return etts.metadata, nil
	}

	return nil, fmt.Errorf("speech generation not finished, first call one of SaveToFile() or GetSpeech() methods")
}

package edgeTTS

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type edgeTTS struct {
	communicator *communicate
	task         *communicateTextTask
	outCome      io.WriteCloser
}

type Args struct {
	Text           string
	Voice          string
	Proxy          string
	Rate           string
	Volume         string
	WriteMedia     string
	WriteSubtitles string
}

func Speak(args Args) error {
	if err := validateArgs(args); err != nil {
		return err
	}

	tts, err := newTTS(args)
	if err != nil {
		return err
	}

	if err := tts.speak(); err != nil {
		return err
	}

	return nil
}

func validateArgs(args Args) error {
	if args.Text == "" {
		return fmt.Errorf("'Args.Text' should contain text to speak but empty string set")
	}

	if args.Voice == "" {
		return fmt.Errorf("'Args.Voice' should contain 'Voice.ShortName' to speak with but empty string set")
	}

	if args.WriteMedia == "" {
		return fmt.Errorf("'Args.WriteMedia' should contain mp3 filename speach should be saved to but empty string set")
	}

	return nil
}

func newTTS(args Args) (*edgeTTS, error) {
	if _, err := os.Stat(args.WriteMedia); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(args.WriteMedia), 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create dir: %v", err)
		}
	}

	tts := newCommunicate().
		withVoice(args.Voice).
		withRate(args.Rate).
		withVolume(args.Volume).
		withProxy(args.Proxy)

	file, err := os.OpenFile(args.WriteMedia, os.O_WRONLY|os.O_APPEND|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	eTTS := &edgeTTS{
		communicator: tts,
		outCome:      file,
		task: &communicateTextTask{
			text: args.Text,
		},
	}

	return eTTS, nil
}

func (eTTS *edgeTTS) speak() error {
	defer eTTS.outCome.Close()

	if err := eTTS.communicator.process(eTTS.task); err != nil {
		return fmt.Errorf("failed to request server: %w", err)
	}

	if _, err := eTTS.outCome.Write(eTTS.task.speechData); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

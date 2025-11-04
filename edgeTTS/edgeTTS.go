package edgeTTS

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type edgeTTS struct {
	communicator    *communicate
	task            *communicateTextTask
	audioOutCome    io.WriteCloser
	writeMeta       bool
	metadataOutCome io.WriteCloser
}

type Args struct {
	Text          string
	Voice         string
	Rate          string
	Volume        string
	WriteMedia    string
	WriteMetadata string
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
	// create directory for audio file
	if _, err := os.Stat(args.WriteMedia); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(args.WriteMedia), 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create dir: %v", err)
		}
	}

	// open audio file
	audioFile, err := os.OpenFile(args.WriteMedia, os.O_WRONLY|os.O_APPEND|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	var metadataFile *os.File
	if args.WriteMetadata != "" {
		// create directory for metadata file
		if _, err := os.Stat(args.WriteMetadata); os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(args.WriteMetadata), 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create dir: %v", err)
			}
		}

		// open metadata file
		metadataFile, err = os.OpenFile(args.WriteMetadata, os.O_WRONLY|os.O_APPEND|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %v", err)
		}
	}

	tts := newCommunicate().
		withVoice(args.Voice).
		withRate(args.Rate).
		withVolume(args.Volume)

	eTTS := &edgeTTS{
		communicator:    tts,
		audioOutCome:    audioFile,
		writeMeta:       metadataFile != nil,
		metadataOutCome: metadataFile,
		task: &communicateTextTask{
			text: args.Text,
		},
	}

	return eTTS, nil
}

func (eTTS *edgeTTS) speak() error {
	if err := eTTS.communicator.process(eTTS.task); err != nil {
		return fmt.Errorf("failed to request server: %v", err)
	}

	var wg sync.WaitGroup
	done := make(chan struct{})
	errors := make(chan error)

	wg.Go(func() {
		defer eTTS.audioOutCome.Close()

		if _, err := eTTS.audioOutCome.Write(eTTS.task.speechData); err != nil {
			errors <- fmt.Errorf("failed to write to audio file: %v", err)
			close(errors)
		}
	})

	if eTTS.writeMeta {
		wg.Go(func() {
			defer eTTS.metadataOutCome.Close()

			buf, err := json.Marshal(eTTS.task.metaData)
			if err != nil {
				errors <- fmt.Errorf("failed to marshal metadata: %v", err)
				close(errors)
				return
			}

			if _, err := eTTS.metadataOutCome.Write(buf); err != nil {
				errors <- fmt.Errorf("failed to write to metadata file: %v", err)
				close(errors)
			}
		})
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errors:
		return err
	case <-done:
		return nil
	}
}

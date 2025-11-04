package main

import (
	"fmt"

	"github.com/kolonist/edge-tts-go-v2/edgeTTS"
)

func main() {
	fmt.Println("Trying to get voices list...")
	voices, err := edgeTTS.ListVoices()
	if err != nil {
		fmt.Println("Error: %w", err)
	}

	voice := ""

	fmt.Println("Voices:")
	for i, v := range voices {
		fmt.Printf(
			"    %d: locale: %s, gender: %s, short name: %s\n",
			i,
			v.Locale,
			v.Gender,
			v.ShortName,
		)

		if voice == "" && v.Locale == "en-US" && v.Gender == "Male" {
			voice = v.ShortName
		}
	}
	fmt.Println("")

	filename := "./sample.mp3"
	text := "edge-tts-go-v2 is a golang module that allows you to use Microsoft Edge's online text-to-speech service from within your golang code or using the provided edge-tts-go-v2 command"
	fmt.Printf(
		"Speak '%s' to audio file '%s' using voice '%s'...\n",
		text,
		filename,
		voice,
	)

	args := edgeTTS.Args{
		Voice:          voice,
		Text:           text,
		Rate:           "+25%",
		WriteMedia:     filename,
		WriteSubtitles: "./subtitles.txt",
	}

	err = edgeTTS.Speak(args)
	if err != nil {
		fmt.Printf("Error trying to convers text to speach:\n%s\n", err.Error())
		return
	}

	fmt.Printf("Success! Listen spoken text in '%s'\n", filename)
}

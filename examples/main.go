package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/kolonist/edgetts"
)

func main() {
	fmt.Println("Trying to get voices list...")
	voices, err := edgetts.ListVoices(context.TODO())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("Voices:")

	voice := ""
	for i, v := range voices {
		// fmt.Printf(
		// 	"    %d: locale: %s, gender: %s, short name: %s\n",
		// 	i,
		// 	v.Locale,
		// 	v.Gender,
		// 	v.ShortName,
		// )

		fmt.Println(i, ":")
		fmt.Println("Name:", v.Name)
		fmt.Println("ShortName:", v.ShortName)
		fmt.Println("Gender:", v.Gender)
		fmt.Println("Locale:", v.Locale)
		fmt.Println("SuggestedCodec:", v.SuggestedCodec)
		fmt.Println("FriendlyName:", v.FriendlyName)
		fmt.Println("Status:", v.Status)
		fmt.Println("oiceTag.ContentCategories:", strings.Join(v.VoiceTag.ContentCategories, ", "))
		fmt.Println("VoiceTag.VoicePersonalities:", strings.Join(v.VoiceTag.VoicePersonalities, ", "))
		fmt.Println("------------------------------------------------------")

		if voice == "" && v.Locale == "en-US" && v.Gender == "Male" {
			voice = v.ShortName
		}
	}

	fmt.Println("")

	filename := "./sample.mp3"
	text := "edgetts is a golang module that allows you to use Microsoft Edge's online text-to-speech service directly from your golang code"
	fmt.Printf(
		"Speak '%s' to audio file '%s' using voice '%s'...\n",
		text,
		filename,
		voice,
	)

	args := edgetts.Args{
		Voice: voice,
		// Text:         text,
		Rate: "+15%",
		// AudioFile:    filename,
		// MetadataFile: "./subtitles.json",
	}

	tts := edgetts.New(args)

	err = tts.
		Speak(text).
		SaveToFile(filename, edgetts.OutputFormatMp3)

	if err != nil {
		fmt.Printf("Error trying to convert text to speach:\n%s\n", err.Error())
		return
	}

	fmt.Printf("Success! Listen speech in '%s'\n", filename)
}

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
		fmt.Println(i, ":")
		fmt.Println("Name:", v.Name)
		fmt.Println("ShortName:", v.ShortName)
		fmt.Println("Gender:", v.Gender)
		fmt.Println("Locale:", v.Locale)
		fmt.Println("SuggestedCodec:", v.SuggestedCodec)
		fmt.Println("FriendlyName:", v.FriendlyName)
		fmt.Println("Status:", v.Status)
		fmt.Println("VoiceTag.ContentCategories:", "[", strings.Join(v.VoiceTag.ContentCategories, ", "), "]")
		fmt.Println("VoiceTag.VoicePersonalities:", "[", strings.Join(v.VoiceTag.VoicePersonalities, ", "), "]")
		fmt.Println()

		// use first found male english voice
		if voice == "" && v.Locale == "en-US" && v.Gender == "Male" {
			voice = v.ShortName
		}
	}

	fmt.Println()

	filename := "./sample.mp3"
	text := "Edge T T S is a golang module that allows you to use Microsoft Edge's online text-to-speech service directly from your golang code"
	fmt.Printf(
		"Speak '%s' to audio file '%s' using voice '%s'...\n",
		text,
		filename,
		voice,
	)

	args := edgetts.Args{
		Voice: voice,
		Rate:  "+15%",
	}

	speaker := edgetts.
		New(args).
		Speak(text)

	if err := speaker.SaveToFile(context.TODO(), filename, edgetts.OutputFormatMp3); err != nil {
		fmt.Printf("Error trying to synthesize speech:\n%s\n", err.Error())
		return
	}

	fmt.Printf("Success! Listen speech in '%s'\n", filename)

	metadata, err := speaker.GetMetadata()
	if err != nil {
		fmt.Printf("Error trying to get metadata of synthesized speech:\n%s\n", err.Error())
	}

	fmt.Println("")
	fmt.Println("Metadata:")
	fmt.Println("[ Offset ms: Word - Duration ms ]")

	for _, word := range metadata {
		fmt.Printf("%5d: %s - %d\n", word.Offset, word.Text, word.Duration)
	}
}

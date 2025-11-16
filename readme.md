# edgetts

`edgetts` is a golang module that allows you to use Microsoft Edge's online text-to-speech service in your golang projects.

## Installation

To install it, run the following command:

    $ go install github.com/kolonist/edgetts

## Usage

```go
package main

import (
	"context"
	"github.com/kolonist/edgetts"
)

func main() {
	args := edgetts.Args{
		// set voice to use in speech synthesys
		Voice: "en-US-AlloyTurboMultilingualNeural",
	}

	// generate sound in mp3 format and save it to the file
	err := edgetts.
		New(args).
		Speak("Text I need to speak now").
		SaveToFile(context.TODO(), "./sample.mp3", edgetts.OutputFormatMp3)
}
```

You can find more complex example in [/examples](https://github.com/kolonist/edgetts/tree/main/examples) folder.

## Thanks

I used the following projects as sources of inspiration:

* https://github.com/rany2/edge-tts (similar library for Python)
* https://github.com/surfaceyu/edge-tts-go (Python library rewritten in Go but not working now)

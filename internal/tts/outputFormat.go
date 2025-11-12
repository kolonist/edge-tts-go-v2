package tts

type OutputFormat int8

const (
	OutputFormatMp3 OutputFormat = iota
	OutputFormatWebm
	OutputFormatOgg
	OutputFormatRaw22050
	OutputFormatRaw44100
)

func (f OutputFormat) String() string {
	switch f {
	case OutputFormatWebm:
		return "webm-24khz-16bit-24kbps-mono-opus"
	case OutputFormatOgg:
		return "ogg-24khz-16bit-mono-opus"
	case OutputFormatRaw22050:
		return "raw-22050hz-16bit-mono-pcm"
	case OutputFormatRaw44100:
		return "raw-44100hz-16bit-mono-pcm"
	case OutputFormatMp3:
		fallthrough
	default:
		return "audio-24khz-48kbitrate-mono-mp3"
	}
}

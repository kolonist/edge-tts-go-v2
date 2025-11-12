package tts

import (
	"fmt"
	"regexp"
)

type Args struct {
	// Voice to speak with
	Voice string

	// Volume delta, e.g. "+10%" or "-20%"
	Volume string

	// Rate delta, e.g. "+10%" or "-20%"
	Rate string
}

// func (args *Args) getText() (string, error) {
// 	if args.Text == "" {
// 		return "", fmt.Errorf("text not specified")
// 	}

// 	return args.Text, nil
// }

func (args *Args) GetVoice() (string, error) {
	if args.Voice == "" {
		return "", fmt.Errorf("voice not specified")
	}

	// "en-US-AvaMultilingualNeural" -> ("en", "US", "AvaMultilingualNeural")
	// "zh-CN-guangxi-YunqiNeural" -> ("zh", "CN-guangxi", "YunqiNeural")
	match := regexp.MustCompile(`^([a-z]{2,})-([a-zA-Z-]{2,})-([^\-]+Neural)$`).FindStringSubmatch(args.Voice)
	if match == nil {
		return "", fmt.Errorf("voice has wrong format, use 'ListVoices()' and 'ShortName' field go get correct Voice value")
	}

	lang := match[1]
	region := match[2]
	name := match[3]

	// "Microsoft Server Speech Text to Speech Voice (en-US, AvaMultilingualNeural)"
	// "Microsoft Server Speech Text to Speech Voice (zh-CN-guangxi, YunqiNeural)"
	voice := "Microsoft Server Speech Text to Speech Voice (" + lang + "-" + region + ", " + name + ")"

	return voice, nil
}

func (args *Args) getRate() (string, error) {
	// default value
	if args.Rate == "" {
		return "+0%", nil
	}

	if !regexp.MustCompile(`^[+-]\d+%$`).MatchString(args.Rate) {
		return "", fmt.Errorf("rate should have format '+12%%' or '-34%%'")
	}

	return args.Rate, nil
}

func (args *Args) getVolume() (string, error) {
	// default value
	if args.Volume == "" {
		return "+0%", nil
	}

	if !regexp.MustCompile(`^[+-]\d+%$`).MatchString(args.Volume) {
		return "", fmt.Errorf("volume should have format '+12%%' or '-34%%'")
	}

	return args.Volume, nil
}

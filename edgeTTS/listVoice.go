package edgeTTS

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Voice struct {
	Name           string `json:"Name"`
	ShortName      string `json:"ShortName"`
	Gender         string `json:"Gender"`
	Locale         string `json:"Locale"`
	SuggestedCodec string `json:"SuggestedCodec"`
	FriendlyName   string `json:"FriendlyName"`
	Status         string `json:"Status"`
	Language       string
	VoiceTag       VoiceTag `json:"VoiceTag"`
}
type VoiceTag struct {
	ContentCategories  []string `json:"ContentCategories"`
	VoicePersonalities []string `json:"VoicePersonalities"`
}

func listVoices() ([]Voice, error) {
	// Send GET request to retrieve the list of voices.
	client := http.Client{}
	req, err := http.NewRequest(
		"GET",
		VOICE_LIST+
			"&Sec-MS-GEC="+generate_sec_ms_gec()+
			"&Sec-MS-GEC-Version="+SEC_MS_GEC_VERSION,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for k, v := range BASE_HEADERS {
		req.Header.Set(k, v)
	}
	for k, v := range VOICE_HEADERS {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response.
	var voices []Voice
	err = json.Unmarshal(body, &voices)
	if err != nil {
		return nil, err
	}

	return voices, nil
}

type VoicesManager struct {
	voices       []Voice
	calledCreate bool
}

func (vm *VoicesManager) create(customVoices []Voice) error {
	vm.voices = customVoices
	if customVoices == nil {
		voices, err := listVoices()
		if err != nil {
			return err
		}
		vm.voices = voices
	}
	for i, voice := range vm.voices {
		locale := voice.Locale
		if locale == "" {
			return errors.New("invalid voice locale")
		}
		language := locale[:2]
		vm.voices[i].Language = language
	}
	vm.calledCreate = true
	return nil
}

func (vm *VoicesManager) find(attributes Voice) []Voice {
	if !vm.calledCreate {
		panic("VoicesManager.find() called before VoicesManager.create()")
	}

	var matchingVoices []Voice
	for _, voice := range vm.voices {
		matched := true
		if attributes.Language != "" && attributes.Language != voice.Language {
			matched = false
		}
		if attributes.Name != "" && attributes.Name != voice.Name {
			matched = false
		}
		if attributes.Gender != "" && attributes.Gender != voice.Gender {
			matched = false
		}
		if attributes.Locale != "" && attributes.Locale != voice.Locale {
			matched = false
		}
		if matched {
			matchingVoices = append(matchingVoices, voice)
		}
	}
	return matchingVoices
}

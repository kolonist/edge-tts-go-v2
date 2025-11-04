package edgeTTS

import (
	"encoding/json"
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

func ListVoices() ([]Voice, error) {
	// Send GET request to retrieve the list of voices.
	client := http.Client{}
	req, err := http.NewRequest(
		"GET",
		voiceListUrl+
			"&Sec-MS-GEC="+generateSecMSGEC()+
			"&Sec-MS-GEC-Version="+secMSGECVersion,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for k, v := range baseHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range voiceHeaders {
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

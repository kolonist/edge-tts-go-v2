package edgeTTS

const (
	chromiumFullVersion  = "140.0.3485.14"
	chromiumMajorVersion = "140"
	secMSGECVersion      = "1-" + chromiumFullVersion
	baseURL              = "api.msedgeservices.com/tts/cognitiveservices"
	trustedClientToken   = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	wssURL               = "wss://" + baseURL + "/websocket/v1?Ocp-Apim-Subscription-Key=" + trustedClientToken
	voiceListUrl         = "https://" + baseURL + "/voices/list?Ocp-Apim-Subscription-Key=" + trustedClientToken
)

const (
	chunkTypeAudio        = "Audio"
	chunkTypeWordBoundary = "WordBoundary"
	chunkTypeSessionEnd   = "SessionEnd"
	chunkTypeEnd          = "ChunkEnd"
)

var (
	baseHeaders map[string]string = map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36" +
			" (KHTML, like Gecko) Chrome/" + chromiumMajorVersion + ".0.0.0 Safari/537.36" +
			" Edg/" + chromiumMajorVersion + ".0.0.0",
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "en-US,en;q=0.9",
	}
	wssHeaders map[string]string = map[string]string{
		"Pragma":                 "no-cache",
		"Cache-Control":          "no-cache",
		"Origin":                 "chrome-extension://jdiccldimpdaibmpdkjnbmckianbfold",
		"Sec-WebSocket-Protocol": "synthesize",
		// "Sec-WebSocket-Version":  "13", // don't need
	}
	voiceHeaders map[string]string = map[string]string{
		"Authority": "speech.platform.bing.com",
		"Sec-CH-UA": `" Not;A Brand";v="99", "Microsoft Edge";v="` + chromiumMajorVersion + `",` +
			` "Chromium";v="` + chromiumMajorVersion + `"`,
		"Sec-CH-UA-Mobile": "?0",
		"Accept":           "*/*",
		"Sec-Fetch-Site":   "none",
		"Sec-Fetch-Mode":   "cors",
		"Sec-Fetch-Dest":   "empty",
	}
)

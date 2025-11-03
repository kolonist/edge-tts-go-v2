package edgeTTS

const (
	CHROMIUM_FULL_VERSION  = "140.0.3485.14"
	CHROMIUM_MAJOR_VERSION = "140"
	SEC_MS_GEC_VERSION     = "1-" + CHROMIUM_FULL_VERSION
	BASE_URL               = "api.msedgeservices.com/tts/cognitiveservices"
	TRUSTED_CLIENT_TOKEN   = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	WSS_URL                = "wss://" + BASE_URL + "/websocket/v1?Ocp-Apim-Subscription-Key=" + TRUSTED_CLIENT_TOKEN
	VOICE_LIST             = "https://" + BASE_URL + "/voices/list?Ocp-Apim-Subscription-Key=" + TRUSTED_CLIENT_TOKEN
)

// Locale
const (
	ZhCN = "zh-CN"
	EnUS = "en-US"
)

const (
	ChunkTypeAudio        = "Audio"
	ChunkTypeWordBoundary = "WordBoundary"
	ChunkTypeSessionEnd   = "SessionEnd"
	ChunkTypeEnd          = "ChunkEnd"
)

var (
	BASE_HEADERS map[string]string = map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36" +
			" (KHTML, like Gecko) Chrome/" + CHROMIUM_MAJOR_VERSION + ".0.0.0 Safari/537.36" +
			" Edg/" + CHROMIUM_MAJOR_VERSION + ".0.0.0",
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "en-US,en;q=0.9",
	}
	WSS_HEADERS map[string]string = map[string]string{
		"Pragma":                 "no-cache",
		"Cache-Control":          "no-cache",
		"Origin":                 "chrome-extension://jdiccldimpdaibmpdkjnbmckianbfold",
		"Sec-WebSocket-Protocol": "synthesize",
		"Sec-WebSocket-Version":  "13",
	}
	VOICE_HEADERS map[string]string = map[string]string{
		"Authority": "speech.platform.bing.com",
		"Sec-CH-UA": "\" Not;A Brand\";v=\"99\", \"Microsoft Edge\";v=\"" + CHROMIUM_MAJOR_VERSION + "\"," +
			" \"Chromium\";v=\"" + CHROMIUM_MAJOR_VERSION + "\"",
		"Sec-CH-UA-Mobile": "?0",
		"Accept":           "*/*",
		"Sec-Fetch-Site":   "none",
		"Sec-Fetch-Mode":   "cors",
		"Sec-Fetch-Dest":   "empty",
	}
)

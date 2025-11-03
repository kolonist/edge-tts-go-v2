package edgeTTS

const (
	BASE_URL             = "api.msedgeservices.com/tts/cognitiveservices"
	TRUSTED_CLIENT_TOKEN = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	WSS_URL              = "wss://" + BASE_URL + "/websocket/v1?Ocp-Apim-Subscription-Key=" + TRUSTED_CLIENT_TOKEN
	VOICE_LIST           = "https://" + BASE_URL + "/voices/list?Ocp-Apim-Subscription-Key=" + TRUSTED_CLIENT_TOKEN
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

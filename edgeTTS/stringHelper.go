package edgeTTS

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func uuidWithOutDashes() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func bytesToString(text interface{}) string {
	var testBytes string
	switch v := text.(type) {
	case string:
		testBytes = v
	case []byte:
		testBytes = string(v)
	default:
		panic("str must be string or []byte")
	}
	return testBytes
}

func mkssml(text interface{}, voice string, rate string, volume string) string {
	textStr := bytesToString(text)
	ssml := fmt.Sprintf("<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='en-US'><voice name='%s'><prosody pitch='+0Hz' rate='%s' volume='%s'>%s</prosody></voice></speak>", voice, rate, volume, textStr)
	return ssml
}

func ssmlHeadersPlusData(requestID string, timestamp string, ssml string) string {
	return fmt.Sprintf(
		"X-RequestId:%s\r\n"+
			"Content-Type:application/ssml+xml\r\n"+
			"X-Timestamp:%sZ\r\n"+
			"Path:ssml\r\n\r\n"+
			"%s",
		requestID, timestamp, ssml)
}

func dateToString() string {
	return time.Now().UTC().Format("Mon Jan 02 2006 15:04:05 GMT-0700 (Coordinated Universal Time)")
}

func getHeadersAndData(data interface{}) (map[string]string, []byte, error) {
	var dataBytes []byte
	switch v := data.(type) {
	case string:
		dataBytes = []byte(v)
	case []byte:
		dataBytes = v
	default:
		return nil, nil, fmt.Errorf("data must be string or []byte")
	}

	headers := make(map[string]string)
	lines := bytes.Split(dataBytes[:bytes.Index(dataBytes, []byte("\r\n\r\n"))], []byte("\r\n"))
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) < 2 {
			continue
		}
		key := string(parts[0])
		value := strings.TrimSpace(string(parts[1]))
		headers[key] = value
	}

	return headers, dataBytes[bytes.Index(dataBytes, []byte("\r\n\r\n"))+4:], nil
}

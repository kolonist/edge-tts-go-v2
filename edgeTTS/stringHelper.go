package edgeTTS

import (
	"bytes"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	colon     []byte = []byte(":")
	eol       []byte = []byte("\r\n")
	doubleEol []byte = []byte("\r\n\r\n")
)

func uuidWithoutDashes() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func mkssml(text string, voice string, rate string, volume string) string {
	return "<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='en-US'>" +
		"<voice name='" + voice + "'>" +
		"<prosody pitch='+0Hz' rate='" + rate + "' volume='" + volume + "'>" + text + "</prosody>" +
		"</voice></speak>"
}

func ssmlHeadersPlusData(requestID string, timestamp string, ssml string) string {
	return "X-RequestId:" + requestID + "\r\n" +
		"Content-Type:application/ssml+xml\r\n" +
		"X-Timestamp:" + timestamp + "Z\r\n" +
		"Path:ssml\r\n\r\n" +
		ssml
}

func getCurrentTime() string {
	return time.
		Now().
		UTC().
		Format("Mon Jan 02 2006 15:04:05 GMT-0700 (Coordinated Universal Time)")
}

func getHeadersAndData(data []byte) (map[string]string, []byte) {
	headers := make(map[string]string)

	lines := bytes.SplitSeq(data[:bytes.Index(data, doubleEol)], eol)

	for line := range lines {
		parts := bytes.SplitN(line, colon, 2)
		if len(parts) < 2 {
			continue
		}

		key := string(parts[0])
		value := strings.TrimSpace(string(parts[1]))

		headers[key] = value
	}

	return headers, data[bytes.Index(data, doubleEol)+4:]
}

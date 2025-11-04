package edgeTTS

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

type turnMetaInnerText struct {
	Text         string `json:"Text"`
	Length       int    `json:"Length"`
	BoundaryType string `json:"BoundaryType"`
}

type turnMetaInnerData struct {
	Offset   int               `json:"Offset"`
	Duration int               `json:"Duration"`
	Text     turnMetaInnerText `json:"text"`
}

type turnMetadata struct {
	Type string            `json:"Type"`
	Data turnMetaInnerData `json:"Data"`
}

type turnMeta struct {
	Metadata []turnMetadata `json:"Metadata"`
}

type metadataChunk struct {
	Offset   int    `json:"offset"`
	Duration int    `json:"duration"`
	Text     string `json:"text"`
}

type communicateChunk struct {
	Type string
	Data []byte

	metadataChunk
}

type communicateTextTask struct {
	text       string
	chunk      chan communicateChunk
	speechData []byte
	metaData   []metadataChunk
}

type communicateTextOption struct {
	voice  string
	rate   string
	volume string
}

type communicate struct {
	option communicateTextOption
}

func newCommunicate() *communicate {
	return &communicate{
		option: communicateTextOption{
			voice:  "Microsoft Server Speech Text to Speech Voice (en-US, AvaMultilingualNeural)",
			rate:   "+0%",
			volume: "+0%",
		},
	}
}

func (c *communicate) withVoice(voice string) *communicate {
	if voice == "" {
		return c
	}
	match := regexp.MustCompile(`^([a-z]{2,})-([A-Z]{2,})-(.+Neural)$`).FindStringSubmatch(voice)
	if match != nil {
		lang := match[1]
		region := match[2]
		name := match[3]
		if i := strings.Index(name, "-"); i != -1 {
			region = region + "-" + name[:i]
			name = name[i+1:]
		}
		voice = fmt.Sprintf("Microsoft Server Speech Text to Speech Voice (%s-%s, %s)", lang, region, name)
		if !isValidVoice(voice) {
			return c
		}
		c.option.voice = voice
	}

	return c
}

func (c *communicate) withRate(rate string) *communicate {
	if isValidRate(rate) {
		c.option.rate = rate
	}
	return c
}

func (c *communicate) withVolume(volume string) *communicate {
	if isValidVolume(volume) {
		c.option.volume = volume
	}
	return c
}

func (c *communicate) openWs() (*websocket.Conn, error) {
	headers := http.Header{}

	for k, v := range baseHeaders {
		headers.Set(k, v)
	}
	for k, v := range wssHeaders {
		headers.Set(k, v)
	}

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(
		wssURL+
			"&ConnectionId="+uuidWithoutDashes()+
			"&Sec-MS-GEC="+generateSecMSGEC()+
			"&Sec-MS-GEC-Version="+secMSGECVersion,
		headers,
	)
	if err != nil {
		return nil, fmt.Errorf("dial error: %v", err)
	}
	return conn, nil
}

func (c *communicate) stream(task *communicateTextTask) (chan communicateChunk, error) {
	task.chunk = make(chan communicateChunk)

	conn, err := c.openWs()
	if err != nil {
		return nil, err
	}

	date := getCurrentTime()

	conn.WriteMessage(
		websocket.TextMessage,
		[]byte(
			"X-Timestamp:"+date+"\r\n"+
				"Content-Type:application/json; charset=utf-8\r\n"+
				"Path:speech.config\r\n\r\n"+
				`{"context":{"synthesis":{"audio":{"metadataoptions":{"sentenceBoundaryEnabled":false,"wordBoundaryEnabled":true},"outputFormat":"audio-24khz-48kbitrate-mono-mp3"}}}}`+"\r\n",
		),
	)
	conn.WriteMessage(
		websocket.TextMessage,
		[]byte(
			ssmlHeadersPlusData(
				uuidWithoutDashes(),
				date,
				mkssml(task.text, c.option.voice, c.option.rate, c.option.volume),
			),
		),
	)

	go func() {
		// download indicates whether we should be expecting audio data,
		// this is so what we avoid getting binary data from the websocket
		// and falsely thinking it's audio data.
		downloadAudio := false

		// audio_was_received indicates whether we have received audio data
		// from the websocket. This is so we can raise an exception if we
		// don't receive any audio data.
		// audioWasReceived := false

		for {
			// read message
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			switch messageType {
			case websocket.TextMessage:
				parameters, data := getHeadersAndData(data)
				path := parameters["Path"]

				if path == "turn.start" {
					downloadAudio = true
				} else if path == "turn.end" {
					downloadAudio = false
					task.chunk <- communicateChunk{
						Type: chunkTypeEnd,
					}
				} else if path == "audio.metadata" {
					meta := &turnMeta{}
					if err := json.Unmarshal(data, meta); err != nil {
						log.Fatalf("We received a text message, but unmarshal failed.")
					}
					for _, v := range meta.Metadata {
						if v.Type == chunkTypeWordBoundary {
							task.chunk <- communicateChunk{
								Type: v.Type,
								metadataChunk: metadataChunk{
									Offset:   metadataDrationToMilliseconds(v.Data.Offset),
									Duration: metadataDrationToMilliseconds(v.Data.Duration),
									Text:     v.Data.Text.Text,
								},
							}
						} else if v.Type == chunkTypeSessionEnd {
							continue
						} else {
							log.Fatalf("Unknown metadata type: %s", v.Type)
						}
					}
				} else if path != "response" {
					log.Fatalf("The response from the service is not recognized.\n%s", data)
				}
			case websocket.BinaryMessage:
				if !downloadAudio {
					log.Fatalf("We received a binary message, but we are not expecting one.")
				}
				if len(data) < 2 {
					log.Fatalf("We received a binary message, but it is missing the header length.")
				}
				headerLength := int(binary.BigEndian.Uint16(data[:2]))
				if len(data) < headerLength+2 {
					log.Fatalf("We received a binary message, but it is missing the audio data.")
				}
				task.chunk <- communicateChunk{
					Type: chunkTypeAudio,
					Data: data[headerLength+2:],
				}
				// audioWasReceived = true
			}
		}
	}()

	return task.chunk, nil
}

func (c *communicate) process(task *communicateTextTask) error {
	chunk, err := c.stream(task)
	if err != nil {
		return err
	}

	for {
		v, ok := <-chunk
		if ok {
			if v.Type == chunkTypeAudio {
				task.speechData = append(task.speechData, v.Data...)
			} else if v.Type == chunkTypeWordBoundary {
				task.metaData = append(
					task.metaData,
					metadataChunk{v.Offset, v.Duration, v.Text},
				)
			} else if v.Type == chunkTypeEnd {
				close(task.chunk)
				break
			}
		}
	}

	return nil
}

func isValidVoice(voice string) bool {
	return regexp.MustCompile(`^Microsoft Server Speech Text to Speech Voice \(.+,.+\)$`).MatchString(voice)
}

func isValidRate(rate string) bool {
	if rate == "" {
		return false
	}
	return regexp.MustCompile(`^[+-]\d+%$`).MatchString(rate)
}

func isValidVolume(volume string) bool {
	if volume == "" {
		return false
	}
	return regexp.MustCompile(`^[+-]\d+%$`).MatchString(volume)
}

func metadataDrationToMilliseconds(time int) int {
	return time / 10_000
}

package edgeTTS

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// {"context": {"serviceTag": "743d56a9126e4649b2af1660975e3520"}}
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

type communicateChunk struct {
	Type     string
	Data     []byte
	Offset   int
	Duration int
	Text     string
}

type CommunicateTextTask struct {
	id     int
	text   string
	option CommunicateTextOption

	chunk      chan communicateChunk
	speechData []byte
}

type CommunicateTextOption struct {
	voice  string
	rate   string
	volume string
}

type Communicate struct {
	option CommunicateTextOption
	proxy  string

	processorLimit int
	tasks          chan *CommunicateTextTask
}

func NewCommunicate() *Communicate {
	return &Communicate{
		option: CommunicateTextOption{
			voice:  "Microsoft Server Speech Text to Speech Voice (zh-CN, XiaoxiaoNeural)",
			rate:   "+0%",
			volume: "+0%",
		},
		processorLimit: 16,
		tasks:          make(chan *CommunicateTextTask, 16),
	}
}

func (c *Communicate) WithVoice(voice string) *Communicate {
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

func (c *Communicate) WithRate(rate string) *Communicate {
	if !isValidRate(rate) {
		return c
	}
	c.option.rate = rate
	return c
}

func (c *Communicate) WithVolume(volume string) *Communicate {
	if !isValidVolume(volume) {
		return c
	}
	c.option.volume = volume
	return c
}

func (c *Communicate) WithProxy(proxy string) *Communicate {
	if proxy == "" {
		return c
	}
	c.proxy = proxy
	return c
}

func (c *Communicate) fillOption(text *CommunicateTextOption) {
	if text.voice == "" {
		text.voice = c.option.voice
	}
	if text.rate == "" {
		text.rate = c.option.rate
	}
	if text.volume == "" {
		text.volume = c.option.volume
	}
}

func (c *Communicate) openWs() *websocket.Conn {
	headers := http.Header{}

	for k, v := range BASE_HEADERS {
		headers.Set(k, v)
	}
	for k, v := range WSS_HEADERS {
		headers.Set(k, v)
	}

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(
		WSS_URL+
			"&ConnectionId="+uuidWithOutDashes()+
			"&Sec-MS-GEC="+generate_sec_ms_gec()+
			"&Sec-MS-GEC-Version="+SEC_MS_GEC_VERSION,
		headers)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return conn
}

func (c *Communicate) close() {

}

func (c *Communicate) stream(text *CommunicateTextTask) chan communicateChunk {
	text.chunk = make(chan communicateChunk)
	// texts := splitTextByByteLength(removeIncompatibleCharacters(c.text), calcMaxMsgSize(c.voice, c.rate, c.volume))
	conn := c.openWs()
	date := dateToString()
	c.fillOption(&text.option)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("X-Timestamp:%s\r\nContent-Type:application/json; charset=utf-8\r\nPath:speech.config\r\n\r\n{\"context\":{\"synthesis\":{\"audio\":{\"metadataoptions\":{\"sentenceBoundaryEnabled\":false,\"wordBoundaryEnabled\":true},\"outputFormat\":\"audio-24khz-48kbitrate-mono-mp3\"}}}}\r\n", date)))
	_ = conn.WriteMessage(websocket.TextMessage, []byte(ssmlHeadersPlusData(uuidWithOutDashes(), date, mkssml(
		text.text, text.option.voice, text.option.rate, text.option.volume,
	))))

	go func() {
		// download indicates whether we should be expecting audio data,
		// this is so what we avoid getting binary data from the websocket
		// and falsely thinking it's audio data.
		downloadAudio := false

		// audio_was_received indicates whether we have received audio data
		// from the websocket. This is so we can raise an exception if we
		// don't receive any audio data.
		// audioWasReceived := false

		// finalUtterance := make(map[int]int)
		for {
			// 读取消息
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			switch messageType {
			case websocket.TextMessage:
				parameters, data, _ := getHeadersAndData(data)
				path := parameters["Path"]
				if path == "turn.start" {
					downloadAudio = true
				} else if path == "turn.end" {
					downloadAudio = false
					text.chunk <- communicateChunk{
						Type: ChunkTypeEnd,
					}
				} else if path == "audio.metadata" {
					meta := &turnMeta{}
					if err := json.Unmarshal(data, meta); err != nil {
						log.Fatalf("We received a text message, but unmarshal failed.")
					}
					for _, v := range meta.Metadata {
						if v.Type == ChunkTypeWordBoundary {
							text.chunk <- communicateChunk{
								Type:     v.Type,
								Offset:   v.Data.Offset,
								Duration: v.Data.Duration,
								Text:     v.Data.Text.Text,
							}
						} else if v.Type == ChunkTypeSessionEnd {
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
				text.chunk <- communicateChunk{
					Type: ChunkTypeAudio,
					Data: data[headerLength+2:],
				}
				// audioWasReceived = true
			}
		}
	}()

	return text.chunk
}

func (c *Communicate) allocateTask(tasks []*CommunicateTextTask) {
	for id, t := range tasks {
		t.id = id
		c.tasks <- t
	}
	close(c.tasks)
}

func (c *Communicate) process(wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range c.tasks {
		chunk := c.stream(t)
		for {
			v, ok := <-chunk
			if ok {
				if v.Type == ChunkTypeAudio {
					t.speechData = append(t.speechData, v.Data...)
					// } else if v.Type == ChunkTypeWordBoundary {
				} else if v.Type == ChunkTypeEnd {
					close(t.chunk)
					break
				}
			}
		}
	}
}

func (c *Communicate) createPool() {
	var wg sync.WaitGroup
	for i := 0; i < c.processorLimit; i++ {
		wg.Add(1)
		go c.process(&wg)
	}
	wg.Wait()
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

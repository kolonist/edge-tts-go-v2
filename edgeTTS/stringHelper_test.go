package edgeTTS

import (
	"testing"
)

func Test_getHeadersAndData(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		want1   []byte
		wantErr bool
	}{
		{
			name: "test-1",
			args: args{
				data: []byte(
					"X-Timestamp:2022-01-01\r\n" +
						"Content-Type:application/json; charset=utf-8\r\n" +
						"Path:speech.config\r\n\r\n" +
						`{"context":{"synthesis":{"audio":{"metadataoptions":{"sentenceBoundaryEnabled":false,"wordBoundaryEnabled":true},"outputFormat":"audio-24khz-48kbitrate-mono-mp3"}}}}`,
				),
			},
			want:    map[string]string{},
			want1:   []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getHeadersAndData(tt.args.data)
			t.Logf("%v \n%v \n", got, got1)
		})
	}
}

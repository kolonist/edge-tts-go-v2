package edgeTTS

import (
	"testing"
)

func Test_listVoices(t *testing.T) {
	tests := []struct {
		name    string
		want    []Voice
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test-1",
			want:    []Voice{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListVoices()
			if len(got) <= 0 {
				t.Errorf("ListVoices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

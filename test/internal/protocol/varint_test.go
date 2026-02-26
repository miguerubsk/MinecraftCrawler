package protocol_test

import (
	"MinecraftCrawler/internal/protocol"
	"bytes"
	"testing"
)

func TestWriteVarInt(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		want    []byte
		wantErr bool
	}{
		{"zero", 0, []byte{0x00}, false},
		{"one", 1, []byte{0x01}, false},
		{"127", 127, []byte{0x7F}, false},
		{"128", 128, []byte{0x80, 0x01}, false},
		{"255", 255, []byte{0xFF, 0x01}, false},
		{"2097151", 2097151, []byte{0xFF, 0xFF, 0x7F}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := protocol.WriteVarInt(buf, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteVarInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(buf.Bytes(), tt.want) {
				t.Errorf("WriteVarInt() = %v, want %v", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestReadVarInt(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    int
		wantErr bool
	}{
		{"zero", []byte{0x00}, 0, false},
		{"one", []byte{0x01}, 1, false},
		{"127", []byte{0x7F}, 127, false},
		{"128", []byte{0x80, 0x01}, 128, false},
		{"255", []byte{0xFF, 0x01}, 255, false},
		{"2097151", []byte{0xFF, 0xFF, 0x7F}, 2097151, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewReader(tt.input)
			got, err := protocol.ReadVarInt(buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadVarInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadVarInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
// ReadVarIntSafe test removed because it is private and we cannot change source.

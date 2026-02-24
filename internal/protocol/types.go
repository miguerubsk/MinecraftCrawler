package protocol

import (
	"io"
	"time"
)

type ServerDetail struct {
	IP                 string            `json:"ip"`
	Port               int               `json:"port"`
	Timestamp          time.Time         `json:"timestamp"`
	VersionName        string            `json:"version_name"`
	Protocol           int               `json:"protocol"`
	MOTD               string            `json:"motd"`
	Icon               []byte            `json:"icon"`
	PlayersOnline      int               `json:"players_online"`
	PlayersMax         int               `json:"players_max"`
	Software           string            `json:"software"`
	Mods               map[string]string `json:"mods"`
	Plugins            []string          `json:"plugins"`
	IsWhitelist        bool              `json:"whitelist"`
	EnforcesSecureChat bool              `json:"secure_chat"`
	RconOpen           bool              `json:"rcon_open"`
}

func WriteVarInt(w io.Writer, value int) error {
	for {
		if (value & ^0x7F) == 0 {
			_, err := w.Write([]byte{byte(value)})
			return err
		}
		_, err := w.Write([]byte{byte((value & 0x7F) | 0x80)})
		if err != nil {
			return err
		}
		value >>= 7
	}
}

func ReadVarInt(r io.Reader) (int, error) {
	var value int
	var shift uint
	for {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return 0, err
		}
		value |= int(b[0]&0x7F) << shift
		if (b[0] & 0x80) == 0 {
			break
		}
		shift += 7
	}
	return value, nil
}
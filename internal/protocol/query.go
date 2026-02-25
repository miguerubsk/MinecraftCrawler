package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

type QueryResult struct {
	Software string            `json:"software"`
	Plugins  []string          `json:"plugins"`
	MapName  string            `json:"map_name"`
	RawKV    map[string]string `json:"raw_kv"`
}

func GetQueryInfo(ip string, port int, timeout time.Duration) (*QueryResult, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("udp", addr, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sessionId := int32(0x01010101 & 0x0F0F0F0F)
	
	handshake := new(bytes.Buffer)
	_, _ = handshake.Write([]byte{0xFE, 0xFD, 0x09})
	_ = binary.Write(handshake, binary.BigEndian, sessionId)
	
	_, _ = conn.Write(handshake.Bytes())
	
	resp := make([]byte, 2048)
	n, err := conn.Read(resp)
	if err != nil || n < 5 {
		return nil, fmt.Errorf("no query response")
	}

	tokenStr := string(resp[5 : n-1])
	var token int32
	if _, err := fmt.Sscanf(tokenStr, "%d", &token); err != nil {
		return nil, err
	}

	statReq := new(bytes.Buffer)
	_, _ = statReq.Write([]byte{0xFE, 0xFD, 0x00})
	_ = binary.Write(statReq, binary.BigEndian, sessionId)
	_ = binary.Write(statReq, binary.BigEndian, token)
	_, _ = statReq.Write([]byte{0x00, 0x00, 0x00, 0x00})

	_, _ = conn.Write(statReq.Bytes())
	n, err = conn.Read(resp)
	if err != nil || n < 11 {
		return nil, fmt.Errorf("stat request failed")
	}

	data := resp[11:n]
	kvData := bytes.Split(data, []byte{0x00, 0x01, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x00, 0x00})
	
	kvPairs := make(map[string]string)
	parts := bytes.Split(kvData[0], []byte{0x00})
	for i := 0; i < len(parts)-1; i += 2 {
		key := string(parts[i])
		if key == "" { break }
		val := string(parts[i+1])
		kvPairs[key] = val
	}

	result := &QueryResult{
		Software: kvPairs["server_mod"],
		MapName:  kvPairs["map"],
		RawKV:    kvPairs,
		Plugins:  []string{},
	}

	if p, ok := kvPairs["plugins"]; ok {
		pParts := strings.Split(p, ":")
		if len(pParts) > 1 {
			rawPlugins := strings.Split(pParts[1], ";")
			for _, pl := range rawPlugins {
				result.Plugins = append(result.Plugins, strings.TrimSpace(pl))
			}
		}
	}

	return result, nil
}
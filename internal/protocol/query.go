package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

// QueryResult contiene la información extraída vía UDP
type QueryResult struct {
	Software string            `json:"software"`
	Plugins  []string          `json:"plugins"`
	MapName  string            `json:"map_name"`
	RawKV    map[string]string `json:"raw_kv"`
}

// GetQueryInfo realiza una consulta UDP Full Stat
func GetQueryInfo(ip string, port int, timeout time.Duration) (*QueryResult, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("udp", addr, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sessionId := int32(0x01010101 & 0x0F0F0F0F)
	
	// 1. Handshake para obtener el Challenge Token
	handshake := new(bytes.Buffer)
	handshake.Write([]byte{0xFE, 0xFD, 0x09}) // Magic + Type
	binary.Write(handshake, binary.BigEndian, sessionId)
	
	conn.Write(handshake.Bytes())
	
	resp := make([]byte, 2048)
	n, err := conn.Read(resp)
	if err != nil || n < 5 {
		return nil, fmt.Errorf("no query response")
	}

	// El token viene como string en los bytes después del byte 5
	tokenStr := string(resp[5 : n-1])
	var token int32
	fmt.Sscanf(tokenStr, "%d", &token)

	// 2. Full Stat Request
	statReq := new(bytes.Buffer)
	statReq.Write([]byte{0xFE, 0xFD, 0x00}) // Type: Stat
	binary.Write(statReq, binary.BigEndian, sessionId)
	binary.Write(statReq, binary.BigEndian, token)
	statReq.Write([]byte{0x00, 0x00, 0x00, 0x00}) // Padding para Full Stat

	conn.Write(statReq.Bytes())
	n, err = conn.Read(resp)
	if err != nil || n < 11 {
		return nil, fmt.Errorf("stat request failed")
	}

	// 3. Parseo de Key-Value pairs
	// Saltamos el header de 11 bytes
	data := resp[11:n]
	
	// El formato es [Key]\x00[Value]\x00... terminando en \x00\x00
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

	// El campo plugins suele venir como: "CraftBukkit on Mac OS X: WorldEdit 7.2.0; Essentials 2.18.1"
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
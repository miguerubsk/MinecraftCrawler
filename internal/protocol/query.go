package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

// QueryFullStat obtiene plugins, mapa y lista de jugadores vía UDP
func QueryFullStat(address string, timeout time.Duration) (map[string]string, []string, error) {
	conn, err := net.DialTimeout("udp", address, timeout)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	// 1. Handshake de Query para obtener un Token
	sessionId := int32(0x01010101 & 0x0F0F0F0F)
	payload := new(bytes.Buffer)
	payload.Write([]byte{0xFE, 0xFD, 0x09}) // Magic + Type (Handshake)
	_ = binary.Write(payload, binary.BigEndian, sessionId)
	
	if _, err := conn.Write(payload.Bytes()); err != nil {
		return nil, nil, err
	}
	
	resp := make([]byte, 1500)
	n, err := conn.Read(resp)
	if err != nil || n < 5 {
		return nil, nil, err
	}
	
	// El token es un string numérico terminado en null
	tokenStr := string(resp[5 : n-1])
	var token int32
	_, _ = fmt.Sscanf(tokenStr, "%d", &token) // Ignoramos error de parseo, token será 0 si falla

	// 2. Full Stat Request
	payload.Reset()
	payload.Write([]byte{0xFE, 0xFD, 0x00}) // Type (Stat)
	_ = binary.Write(payload, binary.BigEndian, sessionId)
	_ = binary.Write(payload, binary.BigEndian, token)
	payload.Write([]byte{0x00, 0x00, 0x00, 0x00}) // Padding para Full Stat

	if _, err := conn.Write(payload.Bytes()); err != nil {
		return nil, nil, err
	}
	n, err = conn.Read(resp)
	if err != nil || n < 11 {
		return nil, nil, err
	}

	// 3. Parsear el lío de bytes (K-V pairs)
	// El formato es: [Padding 11 bytes] [Key] \x00 [Value] \x00 ... \x00\x01player\x00\x00 [Players]
	data := resp[11:n]
	parts := bytes.Split(data, []byte{0x00, 0x01, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x00, 0x00})
	
	kvSection := parts[0]
	kvPairs := make(map[string]string)
	kvSplit := bytes.Split(kvSection, []byte{0x00})
	for i := 0; i < len(kvSplit)-1; i += 2 {
		key := string(kvSplit[i])
		if key == "" { break }
		val := string(kvSplit[i+1])
		kvPairs[key] = val
	}

	plugins := []string{}
	if p, ok := kvPairs["plugins"]; ok {
		// Formato: "Software: Plugin1; Plugin2"
		pParts := strings.Split(p, ":")
		if len(pParts) > 1 {
			plugins = strings.Split(pParts[1], ";")
		}
	}

	return kvPairs, plugins, nil
}
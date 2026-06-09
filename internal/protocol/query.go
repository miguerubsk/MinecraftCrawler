package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type QueryResult struct {
	MOTD          string            `json:"motd"`
	Software      string            `json:"software"`
	Plugins       []string          `json:"plugins"`
	MapName       string            `json:"map_name"`
	PlayersOnline int               `json:"players_online"`
	PlayersMax    int               `json:"players_max"`
	HostPort      int               `json:"host_port"`
	HostName      string            `json:"host_name"`
	RawKV         map[string]string `json:"raw_kv"`
}

// GetQueryInfo attempts to get server information via the Query protocol.
// It tries Full Stat by default and falls back to Basic Stat if it fails.
func GetQueryInfo(ip string, port int, timeout time.Duration) (*QueryResult, error) {
	// Full Stat is the default
	res, err := GetFullQueryInfo(ip, port, timeout)
	if err == nil {
		return res, nil
	}

	// Fallback to Basic Stat
	return GetBasicQueryInfo(ip, port, timeout)
}

func GetHandshakeToken(conn net.Conn, timeout time.Duration) (int32, int32, error) {
	sessionId := int32(0x01010101 & 0x0F0F0F0F)
	
	handshake := new(bytes.Buffer)
	_, _ = handshake.Write([]byte{0xFE, 0xFD, 0x09})
	_ = binary.Write(handshake, binary.BigEndian, sessionId)
	
	_ = conn.SetDeadline(time.Now().Add(timeout))
	_, err := conn.Write(handshake.Bytes())
	if err != nil {
		return 0, 0, err
	}
	
	resp := make([]byte, 128)
	n, err := conn.Read(resp)
	if err != nil || n < 5 {
		return 0, 0, fmt.Errorf("no query response")
	}

	tokenStr := strings.TrimRight(string(resp[5:n]), "\x00")
	var token int32
	if _, err := fmt.Sscanf(tokenStr, "%d", &token); err != nil {
		return 0, 0, err
	}

	return sessionId, token, nil
}

func GetFullQueryInfo(ip string, port int, timeout time.Duration) (*QueryResult, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("udp", addr, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sessionId, token, err := GetHandshakeToken(conn, timeout)
	if err != nil {
		return nil, err
	}

	statReq := new(bytes.Buffer)
	_, _ = statReq.Write([]byte{0xFE, 0xFD, 0x00})
	_ = binary.Write(statReq, binary.BigEndian, sessionId)
	_ = binary.Write(statReq, binary.BigEndian, token)
	_, _ = statReq.Write([]byte{0x00, 0x00, 0x00, 0x00}) // Padding for Full Stat

	_ = conn.SetDeadline(time.Now().Add(timeout))
	_, _ = conn.Write(statReq.Bytes())
	
	resp := make([]byte, 4096)
	n, err := conn.Read(resp)
	if err != nil || n < 16 {
		return nil, fmt.Errorf("full stat request failed")
	}

	// Full stat response structure:
	// [1] type (0x00)
	// [4] session ID
	// [11] padding/magic
	data := resp[16:n]
	
	// Split by double null to separate KV from Player list (if any)
	sections := bytes.Split(data, []byte{0x00, 0x01, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x00, 0x00})
	
	kvPairs := make(map[string]string)
	parts := bytes.Split(sections[0], []byte{0x00})
	for i := 0; i < len(parts)-1; i += 2 {
		key := string(parts[i])
		if key == "" { break }
		val := string(parts[i+1])
		kvPairs[key] = val
	}

	res := &QueryResult{
		MOTD:     kvPairs["hostname"],
		Software: kvPairs["server_mod"],
		MapName:  kvPairs["map"],
		HostName: kvPairs["hostname"],
		RawKV:    kvPairs,
		Plugins:  []string{},
	}

	_, _ = fmt.Sscanf(kvPairs["numplayers"], "%d", &res.PlayersOnline)
	_, _ = fmt.Sscanf(kvPairs["maxplayers"], "%d", &res.PlayersMax)
	_, _ = fmt.Sscanf(kvPairs["hostport"], "%d", &res.HostPort)

	if p, ok := kvPairs["plugins"]; ok {
		pParts := strings.Split(p, ":")
		if len(pParts) > 1 {
			rawPlugins := strings.Split(pParts[1], ";")
			for _, pl := range rawPlugins {
				trimmed := strings.TrimSpace(pl)
				if trimmed != "" {
					res.Plugins = append(res.Plugins, trimmed)
				}
			}
		}
	}

	return res, nil
}

func GetBasicQueryInfo(ip string, port int, timeout time.Duration) (*QueryResult, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("udp", addr, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sessionId, token, err := GetHandshakeToken(conn, timeout)
	if err != nil {
		return nil, err
	}

	statReq := new(bytes.Buffer)
	_, _ = statReq.Write([]byte{0xFE, 0xFD, 0x00})
	_ = binary.Write(statReq, binary.BigEndian, sessionId)
	_ = binary.Write(statReq, binary.BigEndian, token)
	// No padding for Basic Stat

	_ = conn.SetDeadline(time.Now().Add(timeout))
	_, _ = conn.Write(statReq.Bytes())
	
	resp := make([]byte, 1024)
	n, err := conn.Read(resp)
	if err != nil || n < 16 {
		return nil, fmt.Errorf("basic stat request failed")
	}

	reader := bytes.NewReader(resp[5:n])
	
	res := &QueryResult{
		RawKV:   make(map[string]string),
		Plugins: []string{},
	}

	res.MOTD, _ = readNullTerminatedString(reader)
	_, _ = readNullTerminatedString(reader) // gametype
	res.MapName, _ = readNullTerminatedString(reader)
	
	numPlayers, _ := readNullTerminatedString(reader)
	maxPlayers, _ := readNullTerminatedString(reader)
	_, _ = fmt.Sscanf(numPlayers, "%d", &res.PlayersOnline)
	_, _ = fmt.Sscanf(maxPlayers, "%d", &res.PlayersMax)

	var hostPort uint16
	_ = binary.Read(reader, binary.LittleEndian, &hostPort)
	res.HostPort = int(hostPort)
	
	res.HostName, _ = readNullTerminatedString(reader)

	return res, nil
}

func readNullTerminatedString(r io.Reader) (string, error) {
	var buf bytes.Buffer
	b := make([]byte, 1)
	for {
		_, err := r.Read(b)
		if err != nil {
			return buf.String(), err
		}
		if b[0] == 0 {
			break
		}
		buf.Write(b)
	}
	return buf.String(), nil
}

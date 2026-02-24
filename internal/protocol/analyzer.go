package protocol

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type slpResponse struct {
	Version            struct { Name string; Protocol int } `json:"version"`
	Players            struct { Max int; Online int }      `json:"players"`
	Description        interface{}                          `json:"description"`
	Favicon            string                               `json:"favicon"`
	EnforcesSecureChat bool                                 `json:"enforcesSecureChat"`
	ForgeData          struct {
		Mods []struct { ModID string; Version string } `json:"mods"`
	} `json:"forgeData"`
}

func AnalyzeServer(ip string, port int, timeout time.Duration) (*ServerDetail, error) {
	detail := &ServerDetail{IP: ip, Port: port, Timestamp: time.Now(), Mods: make(map[string]string)}

	// 1. Fase TCP: SLP y Whitelist
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// --- Handshake & Status Request ---
	if err := sendHandshake(conn, ip, port, 1); err != nil { return nil, err }
	writePacket(conn, []byte{0x00}) // Status Request

	// Leer respuesta SLP
	_, _ = ReadVarInt(conn) // Longitud
	id, _ := ReadVarInt(conn)
	if id == 0x00 {
		l, _ := ReadVarInt(conn)
		data := make([]byte, l)
		io.ReadFull(conn, data)
		var res slpResponse
		json.Unmarshal(data, &res)
		
		detail.VersionName = res.Version.Name
		detail.Protocol = res.Version.Protocol
		detail.PlayersMax = res.Players.Max
		detail.PlayersOnline = res.Players.Online
		detail.EnforcesSecureChat = res.EnforcesSecureChat
		
		// Procesar Favicon
		if res.Favicon != "" {
			b64 := strings.TrimPrefix(res.Favicon, "data:image/png;base64,")
			img, _ := base64.StdEncoding.DecodeString(b64)
			detail.Icon = img
		}

		// Procesar Mods
		for _, m := range res.ForgeData.Mods {
			detail.Mods[m.ModID] = m.Version
		}
	}

	// 2. Check Whitelist (Nueva conexión para limpiar estado)
	connLogin, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err == nil {
		defer connLogin.Close()
		sendHandshake(connLogin, ip, port, 2) // NextState: Login
		
		// Login Start
		ls := new(bytes.Buffer)
		WriteVarInt(ls, 0x00)
		ls.WriteString("GeminiCrawler") // Dummy name
		writePacket(connLogin, ls.Bytes())

		pLen, _ := ReadVarInt(connLogin)
		if pLen > 0 {
			pID, _ := ReadVarInt(connLogin)
			if pID == 0x00 { // Disconnect en fase Login
				reasonLen, _ := ReadVarInt(connLogin)
				reason := make([]byte, reasonLen)
				io.ReadFull(connLogin, reason)
				if strings.Contains(strings.ToLower(string(reason)), "whitelist") {
					detail.IsWhitelist = true
				}
			}
		}
	}

	return detail, nil
}

func sendHandshake(conn net.Conn, host string, port int, nextState int) error {
	buf := new(bytes.Buffer)
	WriteVarInt(buf, 0x00) // Packet ID
	WriteVarInt(buf, 763)  // Protocol
	WriteVarInt(buf, len(host))
	buf.WriteString(host)
	binary.Write(buf, binary.BigEndian, uint16(port))
	WriteVarInt(buf, nextState)
	return writePacket(conn, buf.Bytes())
}

func writePacket(conn net.Conn, data []byte) error {
	buf := new(bytes.Buffer)
	WriteVarInt(buf, len(data))
	buf.Write(data)
	_, err := conn.Write(buf.Bytes())
	return err
}
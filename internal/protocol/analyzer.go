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

type StatusResponse struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
	} `json:"players"`
	Description        interface{} `json:"description"`
	Favicon            string      `json:"favicon"`
	EnforcesSecureChat bool        `json:"enforcesSecureChat"`
	ForgeData          struct {
		Mods []struct {
			ModID   string `json:"modid"`
			Version string `json:"version"`
		} `json:"mods"`
	} `json:"forgeData"`
}

func sendHandshake(conn net.Conn, host string, port int, protocol int, nextState int) error {
	var buf bytes.Buffer
	_ = WriteVarInt(&buf, 0x00)
	_ = WriteVarInt(&buf, protocol)
	_ = WriteVarInt(&buf, len(host))
	_, _ = buf.WriteString(host)
	_ = binary.Write(&buf, binary.BigEndian, uint16(port))
	_ = WriteVarInt(&buf, nextState)

	frame := new(bytes.Buffer)
	_ = WriteVarInt(frame, buf.Len())
	_, _ = frame.Write(buf.Bytes())
	_, err := conn.Write(frame.Bytes())
	return err
}

func AnalyzeServer(ip string, port int, timeout time.Duration) (*ServerDetail, error) {
	detail := &ServerDetail{
		IP: ip, Port: port, Timestamp: time.Now(), Mods: make(map[string]string),
	}

	if port == 25575 {
		return analyzeRcon(detail, timeout)
	}

	status, err := GetServerStatus(ip, port, timeout)
	if err != nil {
		return nil, err
	}

	detail.VersionName = status.Version.Name
	detail.Protocol = status.Version.Protocol
	detail.PlayersMax = status.Players.Max
	detail.PlayersOnline = status.Players.Online
	detail.EnforcesSecureChat = status.EnforcesSecureChat

	if status.Favicon != "" {
		b64 := strings.TrimPrefix(status.Favicon, "data:image/png;base64,")
		img, _ := base64.StdEncoding.DecodeString(b64)
		detail.Icon = img
	}

	for _, m := range status.ForgeData.Mods {
		detail.Mods[m.ModID] = m.Version
	}

	connLogin, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err == nil {
		defer connLogin.Close()
		_ = sendHandshake(connLogin, ip, port, detail.Protocol, 2)

		ls := new(bytes.Buffer)
		_ = WriteVarInt(ls, 0x00)
		username := "GeminiCrawler"
		_ = WriteVarInt(ls, len(username))
		_, _ = ls.WriteString(username)
		
		uuid := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF}
		if detail.Protocol >= 764 {
			_, _ = ls.Write(uuid)
		} else if detail.Protocol >= 759 {
			_ = ls.WriteByte(0x01)
			_, _ = ls.Write(uuid)
		}

		frame := new(bytes.Buffer)
		_ = WriteVarInt(frame, ls.Len())
		_, _ = frame.Write(ls.Bytes())
		_, err = connLogin.Write(frame.Bytes()) // Error checked
		if err == nil {
			_, err = ReadVarInt(connLogin)
			if err == nil {
				pID, _ := ReadVarInt(connLogin)
				if pID == 0x00 {
					rLen, _ := ReadVarInt(connLogin)
					reason := make([]byte, rLen)
					_, _ = io.ReadFull(connLogin, reason) // Silenced explicitly
					msg := strings.ToLower(string(reason))
					if strings.Contains(msg, "whitelist") || strings.Contains(msg, "not on the list") {
						detail.IsWhitelist = true
					}
				}
			}
		}
	}

	query, err := GetQueryInfo(ip, port, 2*time.Second)
	if err == nil {
		detail.Plugins = query.Plugins
		if query.Software != "" {
			detail.Software = query.Software
		}
	}

	return detail, nil
}

func analyzeRcon(detail *ServerDetail, timeout time.Duration) (*ServerDetail, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", detail.IP, detail.Port), timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	payload := ""
	packet := new(bytes.Buffer)
	_ = binary.Write(packet, binary.LittleEndian, int32(len(payload)+10))
	_ = binary.Write(packet, binary.LittleEndian, int32(1))
	_ = binary.Write(packet, binary.LittleEndian, int32(3))
	_, _ = packet.WriteString(payload)
	_, _ = packet.Write([]byte{0x00, 0x00})

	_ = conn.SetDeadline(time.Now().Add(timeout))
	_, err = conn.Write(packet.Bytes())
	if err != nil {
		return nil, err
	}

	detail.RconOpen = true
	detail.Software = "RCON Service"
	return detail, nil
}

func GetServerStatus(host string, port int, timeout time.Duration) (*StatusResponse, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil { return nil, err }
	defer conn.Close()

	_ = sendHandshake(conn, host, port, 763, 1)
	sr := new(bytes.Buffer)
	_ = WriteVarInt(sr, 0x00)
	f := new(bytes.Buffer)
	_ = WriteVarInt(f, sr.Len())
	_, _ = f.Write(sr.Bytes())
	_, _ = conn.Write(f.Bytes())

	_, err = ReadVarInt(conn)
	if err != nil { return nil, err }
	pID, _ := ReadVarInt(conn)
	if pID != 0x00 { return nil, fmt.Errorf("id error") }
	jLen, err := ReadVarIntSafe(conn)
	if err != nil { return nil, err }
	
	jBytes := make([]byte, jLen)
	_, err = io.ReadFull(conn, jBytes)
	if err != nil { return nil, err }

	var res StatusResponse
	if err := json.Unmarshal(jBytes, &res); err != nil { return nil, err }
	return &res, nil
}

func ReadVarIntSafe(r io.Reader) (int, error) {
	var value int
	var shift uint
	for {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil { return 0, err }
		value |= int(b[0]&0x7F) << shift
		if (b[0] & 0x80) == 0 { break }
		shift += 7
		if shift > 35 { return 0, fmt.Errorf("varint too big") }
	}
	return value, nil
}
package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

// Estructura simplificada del JSON que devuelve el servidor
type StatusResponse struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
	} `json:"players"`
	Description interface{} `json:"description"` // Puede ser string o dict
	Favicon     string      `json:"favicon"`     // Base64 PNG
	ForgeData   struct {
		Mods []struct {
			ModID   string `json:"modid"`
			Version string `json:"version"`
		} `json:"mods"`
	} `json:"forgeData"`
}

func GetServerStatus(host string, port int, timeout time.Duration) (*StatusResponse, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// --- 1. Enviar PAQUETE HANDSHAKE ---
	var buf bytes.Buffer
	WriteVarInt(&buf, 0x00)                 // Packet ID: Handshake
	WriteVarInt(&buf, 763)                  // Protocol Version (1.20.1 por defecto)
	WriteVarInt(&buf, len(host))            // Host Length
	buf.WriteString(host)                   // Host
	binary.Write(&buf, binary.BigEndian, uint16(port))
	WriteVarInt(&buf, 1)                    // Next State: 1 (Status)

	// Paquete final: [Longitud][Datos]
	fullPacket := new(bytes.Buffer)
	WriteVarInt(fullPacket, buf.Len())
	fullPacket.Write(buf.Bytes())
	conn.Write(fullPacket.Bytes())

	// --- 2. Enviar STATUS REQUEST ---
	statusReq := new(bytes.Buffer)
	WriteVarInt(statusReq, 1)    // Longitud del paquete
	WriteVarInt(statusReq, 0x00) // Packet ID: Status Request
	conn.Write(statusReq.Bytes())

	// --- 3. LEER RESPUESTA ---
	_, err = ReadVarInt(conn) 
	if err != nil {
		return nil, err
	}
	packetID, _ := ReadVarInt(conn)
	if packetID != 0x00 {
		return nil, fmt.Errorf("packet ID error")
	}

	jsonLen, _ := ReadVarInt(conn)
	jsonBytes := make([]byte, jsonLen)
	_, err = io.ReadFull(conn, jsonBytes)
	if err != nil {
		return nil, err
	}

	var response StatusResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
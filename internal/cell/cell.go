package torcell

import (
	"encoding/binary"
	"errors"
	"net"
)

func NewCell(circuitID uint16, cellType CellType, payload []byte) (*Cell, error) {
	if len(payload) > CellSize-3 {
		return nil, errors.New("too big payload")
	}

	fullPayload := make([]byte, CellSize-3)
	copy(fullPayload, payload)

	return &Cell{
		CircuitID: circuitID,
		Type:      cellType,
		Payload:   fullPayload,
	}, nil
}

func (c *Cell) Serialize() []byte {
	buf := make([]byte, CellSize)
	binary.BigEndian.PutUint16(buf[0:2], c.CircuitID)
	buf[2] = byte(c.Type)
	copy(buf[3:], c.Payload)
	return buf
}

func CreateVersionsCell(versions []uint16) (*Cell, error) {
	if len(versions) == 0 {
		versions = []uint16{4}
	}

	payload := make([]byte, CellSize-3)

	binary.BigEndian.PutUint16(payload[0:2], uint16(len(versions)))

	for i, ver := range versions {
		binary.BigEndian.PutUint16(payload[2+2*i:4+2*i], ver)
	}

	return NewCell(0, CellTypeVersions, payload)
}

func CreateNetinfoCell(clientIP string, timestamp uint32) (*Cell, error) {
	payload := make([]byte, CellSize-3)

	binary.BigEndian.PutUint32(payload[0:4], timestamp)

	payload[4] = 0x04

	payload[9] = 1

	payload[10] = 0x04
	ip := net.ParseIP(clientIP).To4()
	if ip != nil {
		copy(payload[11:15], ip)
	} else {
		copy(payload[11:15], net.IPv4(127, 0, 0, 1))
	}

	return NewCell(0, CellTypeNetinfo, payload)
}

func CreateCertsCell() (*Cell, error) {
	payload := make([]byte, CellSize-3)
	payload[0] = 0

	return NewCell(0, CellTypeCerts, payload)
}

func CreateAuthChallengeCell() (*Cell, error) {
	payload := make([]byte, CellSize-3)
	binary.BigEndian.PutUint16(payload[0:2], 0)

	return NewCell(0, CellTypeAuthChallenge, payload)
}

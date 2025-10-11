package circuit

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"gonion/internal/utils"
	"net"
)

func buildVersionsCell(versions []uint16) []byte {
	cell := make([]byte, 512)

	cell[0] = 0
	cell[1] = 0
	cell[2] = 0
	cell[3] = 0

	binary.BigEndian.PutUint16(cell[3:5], 7)

	binary.BigEndian.PutUint16(cell[4:7], 4)

	var o int = 0
	for i := 6; i < len(versions); i++ {
		binary.BigEndian.AppendUint16(cell[i:], versions[o])
		o++
	}

	return cell
}

func ConnectGuard(guard utils.Relay) (net.Conn, error) {

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", guard.IPv4, guard.ORPort), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}

	cell := buildVersionsCell([]uint16{4})

	_, err = conn.Write(cell)
	if err != nil {
		return nil, err
	}

	resp := make([]byte, 514)

	n, err := conn.Read(resp)

	return conn, err
}

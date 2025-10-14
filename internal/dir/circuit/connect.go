package circuit

import (
	"crypto/tls"
	"fmt"
	torcell "gonion/internal/cell"
	"gonion/internal/utils"
	"net"
	"time"
)

func ConnectToGuard(guard *utils.Relay) (*tls.Conn, error) {

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", guard.IPv4, guard.ORPort), 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	fmt.Println("TCP connection established")

	conn.SetDeadline(time.Now().Add(30 * time.Second))

	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		},
	})

	if err := tlsConn.Handshake(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to guard node")

	fmt.Printf("Consensus fingerprint: %s\n", guard.Identity)

	versionsCell, err := torcell.CreateVersionsCell([]uint16{4, 5})
	if err != nil {
		return nil, err
	}

	if _, err := tlsConn.Write(versionsCell.Serialize()); err != nil {
		return nil, fmt.Errorf("failed to send versions cell: %v", err)
	}

	fmt.Println("Versions cell sent successfully")

	response := make([]byte, 512)
	n, err := tlsConn.Read(response)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	fmt.Printf("Received %d bytes response: %x\n", n, response[:n])
	fmt.Println("bin response:", response)

	return tlsConn, nil
}

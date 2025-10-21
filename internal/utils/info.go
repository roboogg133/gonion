package utils

import (
	"time"
)

type NtorOnionKeyXCert struct {
	Numbers     map[int]int
	Ed25519Cert []byte
}

type Relay struct {
	Nickname        string
	Identity        string
	IPv4            string
	IPv6            string
	ORPort          int
	DirPort         int
	Flags           map[string]bool
	Version         string
	Propertys       map[string][]string
	Bandwidth       int
	Digest          string
	PublicationDate time.Time
	Rules           map[string][]string

	IdentityEd25519   []byte
	MasterKeyEd25519  []byte
	OnionKey          []byte // RSA
	SingningKey       []byte // RSA
	OnionKeyXCert     []byte
	NtorOnionKeyXCert NtorOnionKeyXCert

	RouterSignatureEd25519 []byte
	RoutserSignature       []byte
}

type Consensus struct {
	NetworkStatusVersion int
	VoteStatus           string
	ValidAfter           time.Time
	Relays               map[string]Relay
	Signatures           []string
}

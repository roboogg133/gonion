package utils

import "time"

type Relay struct {
	Nickname        string
	Identity        string
	IPv4            string
	IPv6            string
	ORPort          int
	DirPort         int
	Flags           []string
	Version         string
	Protocols       map[string]string
	Bandwidth       int
	ExitPolicy      string
	Digest          string
	PublicationDate time.Time
}

type Consensus struct {
	NetworkStatusVersion int
	VoteStatus           string
	ValidAfter           time.Time
	Relays               []Relay
	Signatures           []string
}

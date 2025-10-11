package utils

import "time"

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
}

type Consensus struct {
	NetworkStatusVersion int
	VoteStatus           string
	ValidAfter           time.Time
	Relays               []Relay
	Signatures           []string
}

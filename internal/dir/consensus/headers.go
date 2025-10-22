package consensus

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"errors"
	"fmt"
	"gonion/internal/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func validConsensus(file []byte) bool {
	reader := bytes.NewReader(file)
	scanner := bufio.NewScanner(reader)

	var validAfter, validUntil time.Time

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "valid-after ") {

			timeString := strings.TrimLeft(line, "valid-after ")
			validAfter, _ = time.Parse("2006-01-02 15:04:05", timeString)

		} else if strings.HasPrefix(line, "valid-until ") {

			timeString := strings.TrimLeft(line, "valid-until ")
			validUntil, _ = time.Parse("2006-01-02 15:04:05", timeString)
			break
		}

	}
	now := time.Now().UTC()
	if (now.Equal(validAfter) || now.After(validAfter)) && now.Before(validUntil) {
		return true
	}

	return false
}

func ParseConsensus(file []byte) utils.Consensus {
	reader := bytes.NewReader(file)
	scanner := bufio.NewScanner(reader)

	var consensus utils.Consensus
	consensus.Relays = make(map[string]utils.Relay)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "r ") {
			var relay utils.Relay
			lineSepareted := strings.Fields(line)

			relay.Propertys = make(map[string][]string)
			relay.Rules = make(map[string][]string)
			relay.Flags = make(map[string]bool)
			relay.Nickname = lineSepareted[1]
			relay.Identity = lineSepareted[2]
			relay.Digest = lineSepareted[3]
			timing, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", lineSepareted[4], lineSepareted[5]))
			relay.PublicationDate = timing.UTC()
			relay.IPv4 = lineSepareted[6]
			relay.ORPort, _ = strconv.Atoi(lineSepareted[7])
			relay.DirPort, _ = strconv.Atoi(lineSepareted[8])

			for scanner.Scan() {
				line = scanner.Text()
				if strings.HasPrefix(line, "s ") {
					line = strings.TrimLeft(line, "s ")
					lineSepareted := strings.Fields(line)
					for _, v := range lineSepareted {
						relay.Flags[v] = true
					}
				} else if strings.HasPrefix(line, "a ") {
					relay.IPv6 = strings.Fields(line)[1]
				} else if strings.HasPrefix(line, "w ") {
					relay.Bandwidth, _ = strconv.Atoi(strings.TrimLeft(line, "w Bandwidth="))
				} else if strings.HasPrefix(line, "v ") {
					relay.Version = strings.TrimLeft(line, "v Tor ")
				} else if strings.HasPrefix(line, "pr ") {
					line = strings.TrimLeft(line, "pr ")

					lineSepareted := strings.Fields(line)
					for _, v := range lineSepareted {
						value := strings.Split(v, "=")
						if strings.Contains(value[1], "-") {
							rangeNumbers := strings.Split(value[1], "-")
							num1, _ := strconv.Atoi(rangeNumbers[0])
							num2, _ := strconv.Atoi(rangeNumbers[1])

							for i := num1; i <= num2; i++ {
								relay.Propertys[value[0]] = append(relay.Propertys[value[0]], fmt.Sprintf("%d", i))
							}
						} else if strings.Contains(value[1], ",") {
							relay.Propertys[value[0]] = strings.Split(value[1], ",")
						} else {
							relay.Propertys[value[0]] = []string{value[1]}
						}
					}
				} else if strings.HasPrefix(line, "p ") {
					var ruleType string

					line = strings.TrimLeft(line, "p ")

					if strings.HasPrefix(line, "accept") {
						line = strings.TrimLeft(line, "accept ")
						ruleType = "accept"
					} else {
						line = strings.TrimLeft(line, "reject ")
						ruleType = "reject"
					}

					relay.Rules[ruleType] = strings.Split(line, ",")
					break
				}
			}

			consensus.Relays[relay.IPv4] = relay
		}

	}
	return consensus
}

func ParseKeys(relays map[string]utils.Relay) map[string]utils.Relay {
	file, err := getServerTorAll()
	if err != nil {
		return nil
	}
	os.WriteFile("sdafasdfa", file, 0777)
	fmt.Println("got server tor all")

	reader := bytes.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {

		line := scanner.Text()
		if strings.HasPrefix(line, "router ") {
			fLine := strings.Fields(line)
			ipAddr := fLine[2]

			if _, exists := relays[ipAddr]; !exists {
				continue
			}
			fLine = nil
			for scanner.Scan() {
				line = scanner.Text()
				if strings.HasPrefix(line, "identity-ed25519") {
					scanner.Scan()

					var line string
					var fulltext string

					for {
						scanner.Scan()
						line = scanner.Text()
						if strings.HasPrefix(line, "-----END") {
							break
						}
						fulltext = fmt.Sprintf("%s%s", fulltext, line)
					}
					fmt.Println(ipAddr, fulltext)

					relayCopy := relays[ipAddr]
					relayCopy.IdentityEd25519, err = base64.StdEncoding.DecodeString(fulltext)
					if err != nil {
						log.Println(err)
						return nil
					}
					relays[ipAddr] = relayCopy
				} else if strings.HasPrefix(line, "master-key-ed25519 ") {
					stringFields := strings.Fields(line)

					relayCopy := relays[ipAddr]
					relayCopy.MasterKeyEd25519, err = base64.RawStdEncoding.DecodeString(stringFields[1])
					if err != nil {
						log.Println(err)
						return nil
					}
					relays[ipAddr] = relayCopy
				} else if strings.HasPrefix(line, "onion-key") {
					scanner.Scan()

					var line string
					var fulltext string

					for {
						scanner.Scan()
						line = scanner.Text()
						if strings.HasPrefix(line, "-----END") {
							break
						}
						fulltext = fmt.Sprintf("%s%s", fulltext, line)
					}

					relayCopy := relays[ipAddr]
					relayCopy.OnionKey, err = base64.StdEncoding.DecodeString(fulltext)
					if err != nil {
						log.Println(err)
						return nil
					}
					relays[ipAddr] = relayCopy
				} else if strings.HasPrefix(line, "router-signature") {
					fmt.Println("BREAKING LOOP")
					break
				}
			}

		}
	}

	return relays
}

func getServerTorAll() ([]byte, error) {

	consesusList := []string{
		Tor26Info, Moria1Info, DizumInfo,
		GabelmooInfo, DannenbergInfo, MaatuskaInfo,
		LongclawInfo, BastetInfo, FaravaharInfo,
	}

	for _, v := range consesusList {
		fmt.Printf("Dowloading %s\n", v)

		resp, err := http.Get(v)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			continue
		}
		fmt.Println("answer 200")

		var reader io.Reader
		if resp.Header.Get("Content-Encoding") == "deflate" {
			fmt.Println("deflating")
			read, err := zlib.NewReader(resp.Body)
			if err != nil {
				continue
			}
			defer read.Close()
			reader = read
		} else {
			reader = resp.Body
		}

		conensusBlob, err := io.ReadAll(reader)
		if err != nil {
			fmt.Println(err)
			continue
		}

		return conensusBlob, nil

	}

	return nil, errors.New("failed to fetch all /tor/server/all")
}

package consensus

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"gonion/internal/utils"
	"io"
	"net/http"
	"runtime"
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

	reader := bytes.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {

		line := scanner.Text()
		if strings.HasPrefix(line, "router ") {
			fLine := strings.Fields(line)
			ipAddr := fLine[2]
			fLine = nil
			runtime.GC()
			for scanner.Scan() {
				line = scanner.Text()
				if strings.HasPrefix(line, "identity-ed25519") {
					scanner.Scan()
					scanner.Scan()
					line = scanner.Text()
					fulltext := strings.Join()
				}
			}

		}
	}

	return nil
}

func getServerTorAll() ([]byte, error) {

	consesusList := []string{
		Tor26Info, Moria1Info, DizumInfo,
		GabelmooInfo, DannenbergInfo, MaatuskaInfo,
		LongclawInfo, BastetInfo, FaravaharInfo,
	}

	for _, v := range consesusList {

		resp, err := http.Get(v)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			continue
		}

		var reader io.Reader
		if resp.Header.Get("Content-Encoding") == "deflate" {
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
			continue
		}
		if !validConsensus(conensusBlob) {
			continue
		}

		return conensusBlob, nil

	}

	return nil, errors.New("failed to fetch all /tor/server/all")
}

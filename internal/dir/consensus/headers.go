package consensus

import (
	"bufio"
	"bytes"
	"fmt"
	"gonion/internal/utils"
	"strconv"
	"strings"
	"time"
)

func ValidConsensus(file []byte) bool {
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
	if now.After(validAfter) && now.Before(validUntil) {
		return true
	} else if now.Equal(validAfter) {
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

			consensus.Relays = append(consensus.Relays, relay)
		}

	}
	return consensus
}

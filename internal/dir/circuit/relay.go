package circuit

import (
	"gonion/internal/utils"
	"strconv"
	"strings"
)

func SelectGuardNode(relays []utils.Relay) utils.Relay {

	var best utils.Relay

	for _, relay := range relays {
		if relay.Bandwidth > best.Bandwidth && relay.Flags["Guard"] && relay.Flags["Running"] && relay.Flags["Stable"] {
			best = relay
		}
	}
	return best
}

func SelectMiddleNode(relays []utils.Relay) utils.Relay {

	var best utils.Relay

	for _, relay := range relays {
		if relay.Bandwidth > best.Bandwidth && !relay.Flags["Guard"] && !relay.Flags["Exit"] && relay.Flags["Running"] {
			best = relay
		}
	}
	return best
}

func SelectExitNode(relays []utils.Relay, port int) utils.Relay {

	var best utils.Relay

	for _, relay := range relays {
		if relay.Bandwidth > best.Bandwidth && relay.Flags["Exit"] && relay.Flags["Running"] {
			var ruleType string
			if len(relay.Rules["accept"]) > 0 {
				ruleType = "accept"
			} else {
				ruleType = "reject"
			}

			for _, rule := range relay.Rules[ruleType] {
				if strings.Contains(rule, "-") {
					rangeNumbers := strings.Split(rule, "-")
					num1, _ := strconv.Atoi(rangeNumbers[0])
					num2, _ := strconv.Atoi(rangeNumbers[1])

					for i := num1; i <= num2; i++ {
						if i == port {
							best = relay
						}
					}
				} else if strings.Contains(rule, ",") {
					list := strings.Split(rule, ",")
					for _, v := range list {
						rulePort, _ := strconv.Atoi(v)
						if rulePort == port {
							best = relay
						}
					}
				} else {
					rulePort, _ := strconv.Atoi(rule)
					if rulePort == port {
						best = relay
					}
				}
			}
		}
	}
	return best
}

func GetCircuit(relays []utils.Relay, port int) []utils.Relay {
	var circuit []utils.Relay

	circuit = append(circuit, SelectGuardNode(relays))
	circuit = append(circuit, SelectMiddleNode(relays))
	circuit = append(circuit, SelectExitNode(relays, port))

	return circuit
}

package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ronny/ripcalc/ipv4"
)

func main() {
	network, err := ipv4.ParseCIDR("192.168.0.1/24")
	if err != nil {
		slog.Error("failed to parse CIDR", "error", err)
		os.Exit(1)
	}

	err = network.Calculate()
	if err != nil {
		slog.Error("failed to calculate network", "error", err)
		os.Exit(1)
	}

	fmt.Println(network.FormattedText())
}

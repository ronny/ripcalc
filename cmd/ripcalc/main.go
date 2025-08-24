package main

import (
	"fmt"
	"os"

	"github.com/ronny/ripcalc/ipv4"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		printUsage()
		return fmt.Errorf("no CIDR argument provided")
	}

	cidr := os.Args[1]

	// Handle help requests
	if cidr == "-h" || cidr == "--help" || cidr == "help" {
		printUsage()
		return nil
	}

	network, err := ipv4.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("invalid CIDR notation %q: %w", cidr, err)
	}

	err = network.Calculate()
	if err != nil {
		return fmt.Errorf("failed to calculate network: %w", err)
	}

	fmt.Println(network.FormattedText())

	return nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `ripcalc - IPv4 and IPv6 address calculator

Usage:
  ripcalc <CIDR>

Arguments:
  CIDR    IPv4 address in CIDR notation (e.g., 192.168.0.0/24)

Options:
  -h, --help    Show this help message

Examples:
  ripcalc 192.168.0.0/24
  ripcalc 10.0.0.1/16
  ripcalc 172.16.0.0/12

`)
}

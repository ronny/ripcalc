package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/ronny/ripcalc/ipv4"
	"github.com/ronny/ripcalc/ipv6"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	return runWithArgs(os.Args)
}

func runWithArgs(args []string) error {
	// Create a new FlagSet to avoid global flag conflicts in tests
	fs := flag.NewFlagSet("ripcalc", flag.ContinueOnError)
	
	// Define flags
	var showMask = fs.Bool("ipv6-mask", false, "Show netmask and wildcard for IPv6 (always shown for IPv4)")
	var showBinary = fs.Bool("ipv6-binary", false, "Show binary representation for IPv6 (always shown for IPv4)")
	var help = fs.Bool("help", false, "Show help message")
	fs.BoolVar(help, "h", false, "Show help message (shorthand)")

	// Custom usage function
	fs.Usage = func() {
		printUsage()
	}

	// Parse flags
	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	// Handle help requests
	if *help {
		printUsage()
		return nil
	}

	// Check for CIDR argument
	flagArgs := fs.Args()
	if len(flagArgs) < 1 {
		printUsage()
		return fmt.Errorf("no CIDR argument provided")
	}

	cidr := flagArgs[0]

	// Detect IP version and handle accordingly
	if isIPv6CIDR(cidr) {
		return handleIPv6(cidr, *showMask, *showBinary)
	} else {
		return handleIPv4(cidr)
	}
}

func isIPv6CIDR(cidr string) bool {
	// Parse the CIDR to check if it's IPv6
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}

	// Check if it's IPv6 by looking at the length and format
	// IPv6 addresses have 16 bytes, but we need to distinguish from IPv4-mapped
	if ip.To16() == nil {
		return false
	}

	// If the original string contains ":", it's likely IPv6
	// This handles IPv4-mapped IPv6 addresses correctly
	return strings.Contains(cidr, ":")
}

func handleIPv4(cidr string) error {
	network, err := ipv4.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("invalid IPv4 CIDR notation %q: %w", cidr, err)
	}

	err = network.Calculate()
	if err != nil {
		return fmt.Errorf("failed to calculate IPv4 network: %w", err)
	}

	fmt.Println(network.FormattedText())

	return nil
}

func handleIPv6(cidr string, showMask, showBinary bool) error {
	network, err := ipv6.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("invalid IPv6 CIDR notation %q: %w", cidr, err)
	}

	err = network.Calculate()
	if err != nil {
		return fmt.Errorf("failed to calculate IPv6 network: %w", err)
	}

	if showMask && showBinary {
		fmt.Println(network.FormattedTextWithMask())
	} else if showMask {
		fmt.Println(network.FormattedTextWithMaskNoBinary())
	} else if showBinary {
		fmt.Println(network.FormattedTextWithBinary())
	} else {
		fmt.Println(network.FormattedText())
	}

	return nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `ripcalc - IPv4 and IPv6 address calculator

Usage:
  ripcalc [OPTIONS] <CIDR>

Arguments:
  CIDR    IPv4 or IPv6 address in CIDR notation

Options:
  -h, --help         Show this help message
      --ipv6-mask    Show netmask and wildcard for IPv6 (always shown for IPv4)
      --ipv6-binary  Show binary representation for IPv6 (always shown for IPv4)

Examples:
  IPv4:
    ripcalc 192.168.0.0/24
    ripcalc 10.0.0.1/16
    ripcalc 172.16.0.0/12

  IPv6:
    ripcalc 2001:db8::/64
    ripcalc fe80::/10
    ripcalc ::1/128
    ripcalc --ipv6-mask 2001:db8::/64
    ripcalc --ipv6-binary 2001:db8::/64
    ripcalc --ipv6-mask --ipv6-binary 2001:db8::/64

`)
}

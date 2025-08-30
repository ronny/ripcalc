package ipv6

import (
	"fmt"
	"math/big"
	"net"
	"strings"
)

type addressType int

const (
	addressTypeGlobalUnicast addressType = iota
	addressTypeLinkLocal
	addressTypeUniqueLocal
	addressTypeMulticast
	addressTypeLoopback
	addressTypeUnspecified
	addressTypeDocumentation
	addressType6to4
	addressTypeTeredo
	addressTypeIPv4Mapped
	addressTypeReserved
)

func (at addressType) String() string {
	switch at {
	case addressTypeGlobalUnicast:
		return "Internet Routable"
	case addressTypeLinkLocal:
		return "Auto-configured"
	case addressTypeUniqueLocal:
		return "Private"
	case addressTypeMulticast:
		return "Group Communication"
	case addressTypeLoopback:
		return "Host-only"
	case addressTypeUnspecified:
		return "Default/Undefined"
	case addressTypeDocumentation:
		return "RFC Example"
	case addressType6to4:
		return "IPv4 Transition (Deprecated)"
	case addressTypeTeredo:
		return "NAT Traversal"
	case addressTypeIPv4Mapped:
		return "Embedded IPv4"
	case addressTypeReserved:
		return "Reserved"
	default:
		return "Unknown"
	}
}

type addressRange struct {
	network *net.IPNet
	typ     addressType
	class   string
}

var specialRanges = []addressRange{
	{mustParseCIDR("::1/128"), addressTypeLoopback, "Loopback"},
	{mustParseCIDR("::/128"), addressTypeUnspecified, "Unspecified"},
	{mustParseCIDR("fe80::/10"), addressTypeLinkLocal, "Link-Local Unicast"},
	{mustParseCIDR("fc00::/7"), addressTypeUniqueLocal, "Unique Local Address"},
	{mustParseCIDR("ff00::/8"), addressTypeMulticast, "Multicast"},
	{mustParseCIDR("2001:db8::/32"), addressTypeDocumentation, "Documentation"},
	{mustParseCIDR("2002::/16"), addressType6to4, "6to4"},
	{mustParseCIDR("2001::/32"), addressTypeTeredo, "Teredo"},
	{mustParseCIDR("::ffff:0:0/96"), addressTypeIPv4Mapped, "IPv4-Mapped"},
	{mustParseCIDR("2000::/3"), addressTypeGlobalUnicast, "Global Unicast"},
}

func mustParseCIDR(cidr string) *net.IPNet {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(fmt.Sprintf("invalid CIDR in mustParseCIDR: %s", cidr))
	}

	return network
}

type Network struct {
	Address      net.IP
	PrefixLength int
	Network      net.IP
	HostMin      net.IP
	HostMax      net.IP
	HostCount    *big.Int
	Class        string
	Type         string
}

func ParseCIDR(cidr string) (*Network, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("net.ParseCIDR: %w", err)
	}

	if ip.To16() == nil {
		return nil, fmt.Errorf("%w: not an IPv6 address", ErrInvalidAddress)
	}

	// Ensure it's actually IPv6 and not IPv4
	if ip.To4() != nil {
		return nil, fmt.Errorf("%w: IPv4 address provided, expected IPv6", ErrInvalidAddress)
	}

	prefixLen, _ := ipNet.Mask.Size()

	return &Network{
		Address:      ip.To16(),
		PrefixLength: prefixLen,
	}, nil
}

func (n *Network) String() string {
	return fmt.Sprintf("%s/%d", n.Address, n.PrefixLength)
}

func (n *Network) Calculate() error {
	if n.Address == nil {
		return fmt.Errorf("%w: address is nil", ErrInvalidAddress)
	}

	// Calculate network address
	mask := net.CIDRMask(n.PrefixLength, 128)
	n.Network = n.Address.Mask(mask)

	// Calculate host range (first and last addresses in subnet)
	n.HostMin, n.HostMax = calculateHostRange(n.Network, n.PrefixLength)

	// Calculate host count (for display purposes, though it may be massive)
	n.HostCount = calculateHostCount(n.PrefixLength)

	// Classify the address
	n.Class, n.Type = classifyAddress(n.Address)

	return nil
}

func (n *Network) FormattedText() string {
	// Format addresses (no binary, no mask - clean default format)
	addressCompressed := compressIPv6(n.Address)
	networkStr := fmt.Sprintf("%s/%d", compressIPv6(n.Network), n.PrefixLength)
	
	// For display purposes, limit host count to avoid enormous numbers
	hostCountStr := formatHostCount(n.HostCount, n.PrefixLength)

	separator := calculateSeparatorLength(false)
	
	return fmt.Sprintf(
		""+
			"   Address:\t%-40s\n"+
			"    Prefix:\t%-40s\n"+
			"%s\n"+
			"   Network:\t%-40s\n"+
			"First host:\t%-40s\n"+
			" Last host:\t%-40s\n"+
			"Host count:\t%-40s\t%s, %s",
		addressCompressed,
		fmt.Sprintf("/%d", n.PrefixLength),
		separator,
		networkStr,
		compressIPv6(n.HostMin),
		compressIPv6(n.HostMax),
		hostCountStr, n.Class, n.Type,
	)
}

func (n *Network) FormattedTextWithBinary() string {
	// Format addresses with binary representations
	addressCompressed := compressIPv6(n.Address)
	networkStr := fmt.Sprintf("%s/%d", compressIPv6(n.Network), n.PrefixLength)
	
	// Format binary representations with network/host boundary
	addressBinary := FormatBinaryWithMask(n.Address, n.PrefixLength)
	networkBinary := FormatBinaryWithMask(n.Network, n.PrefixLength)
	hostMinBinary := FormatBinaryWithMask(n.HostMin, n.PrefixLength)
	hostMaxBinary := FormatBinaryWithMask(n.HostMax, n.PrefixLength)

	// For display purposes, limit host count to avoid enormous numbers
	hostCountStr := formatHostCount(n.HostCount, n.PrefixLength)

	separator := calculateSeparatorLength(true)
	
	return fmt.Sprintf(
		""+
			"   Address:\t%-40s\t%s\n"+
			"    Prefix:\t%-40s\n"+
			"%s\n"+
			"   Network:\t%-40s\t%s\n"+
			"First host:\t%-40s\t%s\n"+
			" Last host:\t%-40s\t%s\n"+
			"Host count:\t%-40s\t%s, %s",
		addressCompressed, addressBinary,
		fmt.Sprintf("/%d", n.PrefixLength),
		separator,
		networkStr, networkBinary,
		compressIPv6(n.HostMin), hostMinBinary,
		compressIPv6(n.HostMax), hostMaxBinary,
		hostCountStr, n.Class, n.Type,
	)
}

func (n *Network) FormattedTextWithMask() string {
	// Calculate netmask and wildcard
	netmask := calculateIPv6Netmask(n.PrefixLength)
	wildcard := calculateIPv6Wildcard(n.PrefixLength)
	
	// Format addresses
	addressCompressed := compressIPv6(n.Address)
	networkStr := fmt.Sprintf("%s/%d", compressIPv6(n.Network), n.PrefixLength)
	
	// Format binary representations with network/host boundary
	addressBinary := FormatBinaryWithMask(n.Address, n.PrefixLength)
	netmaskBinary := FormatBinaryWithMask(netmask, n.PrefixLength)
	wildcardBinary := FormatBinaryWithMask(wildcard, n.PrefixLength)
	networkBinary := FormatBinaryWithMask(n.Network, n.PrefixLength)
	hostMinBinary := FormatBinaryWithMask(n.HostMin, n.PrefixLength)
	hostMaxBinary := FormatBinaryWithMask(n.HostMax, n.PrefixLength)

	// For display purposes, limit host count to avoid enormous numbers
	hostCountStr := formatHostCount(n.HostCount, n.PrefixLength)

	separator := calculateSeparatorLength(true)
	
	return fmt.Sprintf(
		""+
			"   Address:\t%-40s\t%s\n"+
			"    Prefix:\t%-40s\n"+
			"   Netmask:\t%-40s\t%s\n"+
			"  Wildcard:\t%-40s\t%s\n"+
			"%s\n"+
			"   Network:\t%-40s\t%s\n"+
			"First host:\t%-40s\t%s\n"+
			" Last host:\t%-40s\t%s\n"+
			"Host count:\t%-40s\t%s, %s",
		addressCompressed, addressBinary,
		fmt.Sprintf("/%d", n.PrefixLength),
		compressIPv6(netmask), netmaskBinary,
		compressIPv6(wildcard), wildcardBinary,
		separator,
		networkStr, networkBinary,
		compressIPv6(n.HostMin), hostMinBinary,
		compressIPv6(n.HostMax), hostMaxBinary,
		hostCountStr, n.Class, n.Type,
	)
}

func (n *Network) FormattedTextWithMaskNoBinary() string {
	// Calculate netmask and wildcard
	netmask := calculateIPv6Netmask(n.PrefixLength)
	wildcard := calculateIPv6Wildcard(n.PrefixLength)
	
	// Format addresses
	addressCompressed := compressIPv6(n.Address)
	networkStr := fmt.Sprintf("%s/%d", compressIPv6(n.Network), n.PrefixLength)

	// For display purposes, limit host count to avoid enormous numbers
	hostCountStr := formatHostCount(n.HostCount, n.PrefixLength)

	return fmt.Sprintf(
		""+
			"   Address:\t%-40s\n"+
			"    Prefix:\t%-40s\n"+
			"   Netmask:\t%-40s\n"+
			"  Wildcard:\t%-40s\n"+
			"----------------------------------------------------------------------------\n"+
			"   Network:\t%-40s\n"+
			"First host:\t%-40s\n"+
			" Last host:\t%-40s\n"+
			"Host count:\t%-40s\t%s, %s",
		addressCompressed,
		fmt.Sprintf("/%d", n.PrefixLength),
		compressIPv6(netmask),
		compressIPv6(wildcard),
		networkStr,
		compressIPv6(n.HostMin),
		compressIPv6(n.HostMax),
		hostCountStr, n.Class, n.Type,
	)
}

func calculateHostRange(network net.IP, prefixLen int) (net.IP, net.IP) {
	// For IPv6, we'll calculate the first and last possible addresses
	hostMin := make(net.IP, 16)
	hostMax := make(net.IP, 16)

	copy(hostMin, network)
	copy(hostMax, network)

	// Calculate the inverse mask
	hostBits := 128 - prefixLen
	if hostBits <= 0 {
		// Single address
		return hostMin, hostMax
	}

	// Set all host bits to 1 in hostMax
	bytesToFill := hostBits / 8
	remainingBits := hostBits % 8

	// Start from the end and work backwards
	byteIndex := 15
	for range bytesToFill {
		hostMax[byteIndex] = 0xFF
		byteIndex--
	}

	if remainingBits > 0 && byteIndex >= 0 {
		mask := byte((1 << remainingBits) - 1)
		hostMax[byteIndex] |= mask
	}

	return hostMin, hostMax
}

func calculateHostCount(prefixLen int) *big.Int {
	hostBits := 128 - prefixLen
	if hostBits <= 0 {
		return big.NewInt(1)
	}

	// 2^hostBits
	result := big.NewInt(1)
	result.Lsh(result, uint(hostBits))

	return result
}

func classifyAddress(ip net.IP) (string, string) {
	// Check special ranges in order of specificity
	for _, r := range specialRanges {
		if r.network.Contains(ip) {
			// Special handling for multicast to include scope
			if r.typ == addressTypeMulticast {
				scope := getMulticastScope(ip)
				return fmt.Sprintf("Multicast %s", scope), r.typ.String()
			}

			return r.class, r.typ.String()
		}
	}

	// Default to reserved if no match
	return "Reserved", addressTypeReserved.String()
}

func getMulticastScope(ip net.IP) string {
	if len(ip) != 16 || ip[0] != 0xff {
		return "Unknown"
	}

	// Scope is in the lower 4 bits of the second byte
	scope := ip[1] & 0x0f

	switch scope {
	case 0x1:
		return "Interface-Local"
	case 0x2:
		return "Link-Local"
	case 0x4:
		return "Admin-Local"
	case 0x5:
		return "Site-Local"
	case 0x8:
		return "Organization-Local"
	case 0xe:
		return "Global"
	default:
		return fmt.Sprintf("Scope-%X", scope)
	}
}

func compressIPv6(ip net.IP) string {
	// Use Go's built-in IPv6 compression
	return ip.String()
}

func formatHostCount(count *big.Int, prefixLen int) string {
	// For very large numbers, show in scientific notation or with units
	if prefixLen >= 120 {
		return count.String()
	} else if prefixLen >= 64 {
		// Show in powers notation for readability
		hostBits := 128 - prefixLen
		return fmt.Sprintf("2^%d", hostBits)
	} else {
		// For smaller prefixes, this is astronomical - just show the power
		hostBits := 128 - prefixLen
		return fmt.Sprintf("2^%d (astronomical)", hostBits)
	}
}

// FormatBinary returns the full 128-bit binary representation of an IPv6 address
func FormatBinary(ip net.IP) string {
	if len(ip) != 16 {
		return ""
	}

	var result strings.Builder

	for i, b := range ip {
		if i > 0 && i%2 == 0 {
			result.WriteString(":")
		}

		result.WriteString(fmt.Sprintf("%08b", b))
	}

	return result.String()
}

// FormatBinaryWithMask returns IPv6 binary representation with network/host boundary
func FormatBinaryWithMask(ip net.IP, prefixLength int) string {
	if len(ip) != 16 {
		return ""
	}

	if prefixLength >= 128 || prefixLength <= 0 {
		// No mask division needed
		return FormatBinary(ip)
	}

	var result strings.Builder

	bitCount := 0

	for i, b := range ip {
		// Add colon separator every 2 bytes (16 bits), but not at the start
		if i > 0 && i%2 == 0 {
			result.WriteString(":")
		}

		// Process each bit in the byte
		for bit := 7; bit >= 0; bit-- {
			// Add space at network/host boundary
			if bitCount == prefixLength {
				result.WriteString(" ")
			}

			// Add the bit
			if (b>>bit)&1 == 1 {
				result.WriteString("1")
			} else {
				result.WriteString("0")
			}

			bitCount++
		}
	}

	return result.String()
}

// calculateIPv6Netmask returns the IPv6 netmask for a given prefix length
func calculateIPv6Netmask(prefixLen int) net.IP {
	mask := net.CIDRMask(prefixLen, 128)
	return net.IP(mask)
}

// calculateIPv6Wildcard returns the IPv6 wildcard (inverse mask) for a given prefix length  
func calculateIPv6Wildcard(prefixLen int) net.IP {
	mask := net.CIDRMask(prefixLen, 128)
	wildcard := make(net.IP, 16)
	for i := range wildcard {
		wildcard[i] = ^mask[i]
	}
	return wildcard
}

// calculateSeparatorLength determines the appropriate separator line length based on content
func calculateSeparatorLength(hasBinary bool) string {
	if hasBinary {
		// With binary, the line is much longer: 
		// "   Address:\t" + 40 chars + "\t" + 128-bit binary (roughly 145 chars)
		// Total roughly 200+ characters
		return strings.Repeat("-", 200)
	} else {
		// Without binary, just the address part:
		// "   Address:\t" + 40 chars = roughly 50 characters
		return strings.Repeat("-", 76)
	}
}

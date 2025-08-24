package ipv4

import (
	"fmt"
	"net"
)

type addressType int

const (
	addressTypePublic addressType = iota
	addressTypePrivate
	addressTypeSharedAddressSpace
	addressTypeLinkLocal
	addressTypeLoopback
	addressTypeMulticast
)

func (at addressType) String() string {
	switch at {
	case addressTypePublic:
		return "Public Internet"
	case addressTypePrivate:
		return "Private Internet"
	case addressTypeSharedAddressSpace:
		return "Shared Address Space"
	case addressTypeLinkLocal:
		return "Link Local"
	case addressTypeLoopback:
		return "Loopback"
	case addressTypeMulticast:
		return "Multicast"
	default:
		return "Unknown"
	}
}

type addressRange struct {
	network *net.IPNet
	typ     addressType
}

var specialRanges = []addressRange{
	{mustParseCIDR("192.168.0.0/16"), addressTypePrivate},
	{mustParseCIDR("172.16.0.0/12"), addressTypePrivate},
	{mustParseCIDR("10.0.0.0/8"), addressTypePrivate},
	{mustParseCIDR("100.64.0.0/10"), addressTypeSharedAddressSpace},
	{mustParseCIDR("169.254.0.0/16"), addressTypeLinkLocal},
	{mustParseCIDR("127.0.0.0/8"), addressTypeLoopback},
	{mustParseCIDR("224.0.0.0/4"), addressTypeMulticast},
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
	Netmask      net.IPMask
	Wildcard     net.IP
	Network      net.IP
	Broadcast    net.IP
	HostMin      net.IP
	HostMax      net.IP
	HostCount    uint32
	Class        string
	Type         string
}

func ParseCIDR(cidr string) (*Network, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("net.ParseCIDR: %w", err)
	}

	if ip.To4() == nil {
		return nil, fmt.Errorf("%w: not an IPv4 address", ErrInvalidAddress)
	}

	prefixLen, _ := ipNet.Mask.Size()

	return &Network{
		Address:      ip.To4(),
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

	n.Netmask = net.CIDRMask(n.PrefixLength, 32)
	n.Wildcard = invertMask(net.IP(n.Netmask))
	n.Network = n.Address.Mask(n.Netmask)
	n.Broadcast = calculateBroadcast(n.Network, n.Wildcard)
	n.HostMin, n.HostMax = calculateHostRange(n.Network, n.Broadcast)
	n.HostCount = calculateHostCount(n.PrefixLength)
	n.Class = classifyAddress(n.Address)
	n.Type = classifyAddressType(n.Address).String()

	return nil
}

func (n *Network) FormattedText() string {
	addressBinary := FormatBinaryWithMask(n.Address, n.PrefixLength)
	netmaskBinary := FormatBinaryWithMask(net.IP(n.Netmask), n.PrefixLength)
	wildcardBinary := FormatBinaryWithMask(n.Wildcard, n.PrefixLength)
	networkBinary := FormatBinaryWithMask(n.Network, n.PrefixLength)
	hostMinBinary := FormatBinaryWithMask(n.HostMin, n.PrefixLength)
	hostMaxBinary := FormatBinaryWithMask(n.HostMax, n.PrefixLength)
	broadcastBinary := FormatBinaryWithMask(n.Broadcast, n.PrefixLength)

	typeStr := n.Type

	return fmt.Sprintf("Address:\t%-20s\t%s\n"+
		"Netmask:\t%-15s = %-2d\t%s\n"+
		"Wildcard:\t%-20s\t%s\n"+
		"----------------------------------------------------------------------------\n"+
		"Network:\t%-20s\t%s\n"+
		"HostMin:\t%-20s\t%s\n"+
		"HostMax:\t%-20s\t%s\n"+
		"Broadcast:\t%-20s\t%s\n"+
		"Hosts:\t\t%-20d\tClass %s, %s",
		n.Address.String(), addressBinary,
		net.IP(n.Netmask).String(), n.PrefixLength, netmaskBinary,
		n.Wildcard.String(), wildcardBinary,
		fmt.Sprintf("%s/%d", n.Network.String(), n.PrefixLength), networkBinary,
		n.HostMin.String(), hostMinBinary,
		n.HostMax.String(), hostMaxBinary,
		n.Broadcast.String(), broadcastBinary,
		n.HostCount, n.Class, typeStr,
	)
}

func invertMask(mask net.IP) net.IP {
	wildcard := make(net.IP, 4)
	for i := range 4 {
		wildcard[i] = ^mask[i]
	}

	return wildcard
}

func calculateBroadcast(network, wildcard net.IP) net.IP {
	broadcast := make(net.IP, 4)
	for i := range 4 {
		broadcast[i] = network[i] | wildcard[i]
	}

	return broadcast
}

func calculateHostRange(network, broadcast net.IP) (net.IP, net.IP) {
	hostMin := make(net.IP, 4)
	hostMax := make(net.IP, 4)

	copy(hostMin, network)
	copy(hostMax, broadcast)

	// Host min is network + 1
	hostMin[3]++

	// Host max is broadcast - 1
	hostMax[3]--

	return hostMin, hostMax
}

func calculateHostCount(prefixLen int) uint32 {
	hostBits := 32 - prefixLen
	if hostBits <= 1 {
		return 0
	}

	return (1 << hostBits) - 2 // -2 for network and broadcast
}

func classifyAddress(ip net.IP) string {
	firstOctet := ip[0]
	switch {
	case firstOctet <= 127:
		return "A"
	case firstOctet <= 191:
		return "B"
	case firstOctet <= 223:
		return "C"
	case firstOctet <= 239:
		return "D"
	default:
		return "E"
	}
}

func classifyAddressType(ip net.IP) addressType {
	for _, r := range specialRanges {
		if r.network.Contains(ip) {
			return r.typ
		}
	}

	return addressTypePublic
}

func FormatBinary(ip net.IP) string {
	if len(ip) != 4 {
		return ""
	}

	return fmt.Sprintf("%08b.%08b.%08b.%08b", ip[0], ip[1], ip[2], ip[3])
}

func FormatBinaryWithMask(ip net.IP, prefixLength int) string {
	if len(ip) != 4 {
		return ""
	}

	binary := fmt.Sprintf("%08b%08b%08b%08b", ip[0], ip[1], ip[2], ip[3])

	if prefixLength >= 32 || prefixLength <= 0 {
		// Add dots every 8 bits for readability
		return fmt.Sprintf("%s.%s.%s.%s", binary[0:8], binary[8:16], binary[16:24], binary[24:32])
	}

	// Build result with dots and space at network/host boundary
	result := ""
	for i := 0; i < 32; i++ {
		// Add space at network/host boundary
		if i == prefixLength {
			result += " "
		}

		// Add the bit
		result += string(binary[i])

		// Add dot after every 8th bit, but not at the end
		if (i+1)%8 == 0 && i < 31 {
			result += "."
		}
	}

	return result
}
